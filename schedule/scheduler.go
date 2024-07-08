package schedule

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
	"github.com/robfig/cron/v3"
)

func CallSchedule(c *cron.Cron) {
	// set location
	loc, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		log.Println(err)
	}

	cron.WithLocation(loc)

	schedule := &models.Schedule{}
	schedules := schedule.WithQueryAll()

	scheduleMap := make(map[int]cron.EntryID)
	AddSchedules(c, schedules, scheduleMap)

	// Check for new schedules every minute
	c.AddFunc("@every 1m", func() {
		newSchedules := schedule.WithQueryAll()
		AddSchedules(c, newSchedules, scheduleMap)
	})
}

func AddSchedules(c *cron.Cron, schedules []*models.Schedule, scheduleMap map[int]cron.EntryID) {
	triggered := &models.Triggered{}
	scheduleLog := &models.ScheduleLog{}
	for _, schedule := range schedules {
		if !schedule.Active || schedule.Request == nil {
			continue
		}
		if _, exists := scheduleMap[schedule.ID]; !exists {
			id, err := c.AddFunc(schedule.Timing, func() {
				defer func() {
					if r := recover(); r != nil {
						config.App().Log.Warn("Recovered from panic in schedule", fmt.Sprintf("%v", r))
					}
				}()

				startAt := time.Now()

				client := &http.Client{
					Timeout: time.Duration(schedule.Timeout) * time.Second,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}
				req, err := http.NewRequest(schedule.Request.Method, schedule.Request.Url, strings.NewReader(string(schedule.Request.Content)))
				if err != nil {
					config.App().Log.Warn("Schedule Request Error", err.Error())
					return
				}

				for _, header := range schedule.Request.RequestHeaders {
					req.Header.Set(header.Key, header.Value)
				}

				triggered.Create(schedule.ID)
				scheduleUpdate(schedule, true)
				var resp *http.Response
				for retries := 0; retries < schedule.Retries; retries++ {
					resp, err = client.Do(req)
					if err == nil {
						break
					}
					config.App().Log.Warn("Schedule Do Error, retrying", fmt.Sprintf("Attempt %d/%d: %v", retries+1, schedule.Retries, err.Error()))
					time.Sleep(1 * time.Second)
				}
				scheduleUpdate(schedule, false)
				triggered.Delete(schedule.ID)
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					config.App().Log.Warn("Schedule Body Error", err.Error())
					return
				}
				notification(schedule, body)

				finishAt := time.Now()
				scheduleLog.StartedAt = &startAt
				scheduleLog.FinishedAt = &finishAt
				scheduleLog.Took = float32(finishAt.Sub(startAt).Seconds())
				scheduleLog.Result = string(body)
				scheduleLog.Create(schedule.ID)

				webhooks(schedule)
			})
			if err != nil {
				config.App().Log.Warn("Schedule Error", err.Error())
			} else {
				scheduleMap[schedule.ID] = id
			}
		}
	}
}

func scheduleUpdate(schedule *models.Schedule, running bool) {
	query := "UPDATE schedules SET running=$1 WHERE id=$2"
	err := schedule.Update(query, []any{running, schedule.ID})
	if err != nil {
		config.App().Log.Warn("Schedule Update Error", err.Error())
	}
}

func notification(schedule *models.Schedule, body []byte) {
	if schedule.Notification == nil {
		return
	}
	if schedule.Notification.IsMail {
		for _, mail := range schedule.Notification.NotifyEmails {
			err := config.App().Mail.SetSubject(schedule.Timing + " is running").SetContent(string(body)).SetTo(mail.Email).SendText()
			if err != nil {
				config.App().Log.Warn("Schedule Mail Error", err.Error())
			}
		}
	}
	/* if schedule.Notification.IsMessage {
		for _, message := range schedule.Notification.NotifyMessages {
			// TODO send message
		}
	} */
}

func webhooks(schedule *models.Schedule) {
	for _, webhook := range schedule.Webhooks {
		if webhook.Request == nil {
			continue
		}
		go func() {
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			req, err := http.NewRequest(webhook.Request.Method, webhook.Request.Url, strings.NewReader(string(webhook.Request.Content)))
			if err != nil {
				config.App().Log.Warn("Schedule Webhook Error", err.Error())
				return
			}

			_, err = client.Do(req)
			if err != nil {
				config.App().Log.Warn("Schedule Webhook Error", err.Error())
				return
			}
		}()
	}
}
