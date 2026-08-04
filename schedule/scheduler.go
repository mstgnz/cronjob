package schedule

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/pkg/config"
	"github.com/mstgnz/cronjob/pkg/conn"
	"github.com/mstgnz/cronjob/pkg/logger"
	"github.com/robfig/cron/v3"
)

// entry is what we remember about a schedule already registered on the scheduler.
// Timing is kept so a retimed schedule can be detected and re-registered.
type entry struct {
	ID     cron.EntryID
	Timing string
}

// registry tracks registered entries by schedule id. It is written from the
// reconcile tick, which runs on a cron goroutine, so every access is guarded.
type registry struct {
	mu      sync.Mutex
	entries map[int]entry
}

func newRegistry() *registry {
	return &registry{entries: make(map[int]entry)}
}

func (r *registry) get(scheduleID int) (entry, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[scheduleID]
	return e, ok
}

func (r *registry) set(scheduleID int, e entry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[scheduleID] = e
}

func (r *registry) delete(scheduleID int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, scheduleID)
}

func (r *registry) snapshot() map[int]entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[int]entry, len(r.entries))
	for k, v := range r.entries {
		out[k] = v
	}
	return out
}

// CallSchedule performs the first sync and then reconciles every minute.
// The scheduler itself (location and job chain) is built in pkg/config.
func CallSchedule(c *cron.Cron) {
	schedule := &models.Schedule{}
	reg := newRegistry()

	SyncSchedules(c, schedule.WithQueryAll(), reg)

	if _, err := c.AddFunc("@every 1m", func() {
		SyncSchedules(c, schedule.WithQueryAll(), reg)
	}); err != nil {
		logger.Warn("Schedule Sync Error", err.Error())
	}
}

// SyncSchedules reconciles the scheduler with the active schedules in the database:
// new ones are registered, retimed ones are re-registered, and ones that were
// deactivated or deleted are removed. Without the removal side, disabling a
// schedule in the UI would have no effect until the process restarts.
func SyncSchedules(c *cron.Cron, schedules []*models.Schedule, reg *registry) {
	seen := make(map[int]struct{}, len(schedules))

	for _, schedule := range schedules {
		if !schedule.Active || schedule.Request == nil {
			continue
		}
		seen[schedule.ID] = struct{}{}

		if existing, ok := reg.get(schedule.ID); ok {
			if existing.Timing == schedule.Timing {
				continue
			}
			// timing changed: drop the old entry before registering the new one
			c.Remove(existing.ID)
			reg.delete(schedule.ID)
		}

		currentSchedule := schedule
		id, err := c.AddFunc(currentSchedule.Timing, func() {
			runSchedule(currentSchedule)
		})
		if err != nil {
			logger.Warn("Schedule Error", err.Error())
			continue
		}
		reg.set(currentSchedule.ID, entry{ID: id, Timing: currentSchedule.Timing})
	}

	// anything registered but no longer active in the database has to go
	for scheduleID, e := range reg.snapshot() {
		if _, ok := seen[scheduleID]; ok {
			continue
		}
		c.Remove(e.ID)
		reg.delete(scheduleID)
	}
}

// leaseMargin covers the work around the request itself: notification mail,
// writing the log row and firing the webhooks.
const leaseMargin = 2 * time.Minute

// leaseFor is how long the lock is held before another instance may take it over.
// It is derived from the worst case this schedule can take, all attempts plus the
// pause between them, so a genuinely slow job is never treated as a crashed one.
// The lease is a backstop for a dead instance, not a timeout: a healthy run always
// releases the lock itself when it finishes.
func leaseFor(schedule *models.Schedule) time.Duration {
	timeout := time.Duration(schedule.Timeout) * time.Second
	if timeout <= 0 {
		timeout = conn.DefaultOutboundTimeout
	}

	attempts := schedule.Retries + 1
	worstCase := time.Duration(attempts)*timeout + time.Duration(attempts-1)*time.Second

	return worstCase + leaseMargin
}

// runSchedule executes a single schedule: it fires the configured request,
// notifies, logs the result and triggers the webhooks.
//
// Every replica reaches this function at the same moment for the same schedule, so
// the first thing it does is claim the lock. Only the instance that wins carries on;
// the others return without touching the target endpoint.
func runSchedule(currentSchedule *models.Schedule) {
	// no new work once shutdown started
	if config.IsShutting() {
		return
	}

	config.IncrementRunning()
	defer config.DecrementRunning()

	triggered := &models.Triggered{}
	scheduleLog := &models.ScheduleLog{}
	instanceID := InstanceID()

	acquired, err := triggered.Acquire(currentSchedule.ID, instanceID, leaseFor(currentSchedule))
	if err != nil {
		// The lock is the only thing keeping replicas apart, so a lock that cannot be
		// evaluated means the run has to be skipped rather than risked.
		logger.Warn("Schedule Lock Error", fmt.Sprintf("schedule %d: %v", currentSchedule.ID, err))
		return
	}
	if !acquired {
		// another instance is running this schedule, or this one is still running it
		return
	}
	defer func() {
		if err := triggered.Release(currentSchedule.ID, instanceID); err != nil {
			logger.Warn("Schedule Unlock Error", fmt.Sprintf("schedule %d: %v", currentSchedule.ID, err))
		}
	}()

	startAt := time.Now()

	client := conn.NewOutboundClient(time.Duration(currentSchedule.Timeout) * time.Second)

	scheduleUpdate(currentSchedule, true)
	defer scheduleUpdate(currentSchedule, false)

	var resp *http.Response
	// Retries counts the extra attempts after the first one, so a schedule with
	// retries=0 (the column default) still performs exactly one request.
	attempts := currentSchedule.Retries + 1
	for attempt := 1; attempt <= attempts; attempt++ {
		// the body reader is consumed by a failed attempt, so the request is rebuilt each time
		req, reqErr := http.NewRequest(currentSchedule.Request.Method, currentSchedule.Request.Url, strings.NewReader(string(currentSchedule.Request.Content)))
		if reqErr != nil {
			logger.Warn("Schedule Request Error", reqErr.Error())
			break
		}
		for _, header := range currentSchedule.Request.RequestHeaders {
			req.Header.Set(header.Key, header.Value)
		}

		resp, err = client.Do(req)
		if err == nil {
			break
		}
		logger.Warn("Schedule Do Error, retrying", fmt.Sprintf("Attempt %d/%d: %v", attempt, attempts, err.Error()))
		if attempt < attempts {
			time.Sleep(1 * time.Second)
		}
	}

	// Check if resp is nil before accessing it
	if resp == nil {
		logger.Warn("Schedule Error", "All retries failed, no response received")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn("Schedule Body Error", err.Error())
		return
	}
	notification(currentSchedule, body)

	finishAt := time.Now()
	scheduleLog.StartedAt = &startAt
	scheduleLog.FinishedAt = &finishAt
	scheduleLog.Took = float32(finishAt.Sub(startAt).Seconds())
	scheduleLog.Result = string(body)
	if err := scheduleLog.Create(currentSchedule.ID); err != nil {
		logger.Warn("Schedule Log Error", err.Error())
	}

	webhooks(currentSchedule)
}

func scheduleUpdate(schedule *models.Schedule, running bool) {
	query := "UPDATE schedules SET running=$1 WHERE id=$2"
	err := schedule.Update(query, []any{running, schedule.ID})
	if err != nil {
		logger.Warn("Schedule Update Error", err.Error())
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
				logger.Warn("Schedule Mail Error", err.Error())
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
	var wg sync.WaitGroup
	for _, webhook := range schedule.Webhooks {
		if webhook.Request == nil {
			continue
		}
		currentWebhook := webhook
		wg.Add(1)
		go func() {
			defer wg.Done()
			// this goroutine is outside the cron job chain, so it needs its own recover
			defer func() {
				if r := recover(); r != nil {
					logger.Warn("Recovered from panic in webhook", fmt.Sprintf("%v", r))
				}
			}()

			client := conn.NewOutboundClient(time.Duration(schedule.Timeout) * time.Second)

			req, err := http.NewRequest(currentWebhook.Request.Method, currentWebhook.Request.Url, strings.NewReader(string(currentWebhook.Request.Content)))
			if err != nil {
				logger.Warn("Schedule Webhook Error", err.Error())
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				logger.Warn("Schedule Webhook Error", err.Error())
				return
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}()
	}
	// waiting keeps the webhooks inside the job's lifetime, so graceful shutdown covers them
	wg.Wait()
}
