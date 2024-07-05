package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type Webhook struct {
	ID         int        `json:"id"`
	ScheduleID int        `json:"schedule_id" validate:"required,number"`
	RequestID  int        `json:"request_id" validate:"required,number"`
	Active     bool       `json:"active" validate:"boolean"`
	Schedule   *Schedule  `json:"schedule,omitempty"`
	Request    *Request   `json:"request,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type WebhookUpdate struct {
	ScheduleID int   `json:"schedule_id" validate:"omitempty,number"`
	RequestID  int   `json:"request_id" validate:"omitempty,number"`
	Active     *bool `json:"active" validate:"omitnil,boolean"`
}

func (m *Webhook) Count(userID int) int {
	rowCount := 0

	// prepare count
	stmt, err := config.App().DB.Prepare(config.App().QUERY["WEBHOOKS_COUNT"])
	if err != nil {
		return rowCount
	}

	// query
	rows, err := stmt.Query(userID)
	if err != nil {
		return rowCount
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&rowCount); err != nil {
			return rowCount
		}
	}

	return rowCount
}

func (m *Webhook) Get(id, schedule_id, request_id, user_id int) ([]*Webhook, error) {

	query := strings.TrimSuffix(config.App().QUERY["WEBHOOKS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND w.id=%d", id)
	}

	if schedule_id > 0 {
		query += fmt.Sprintf(" AND w.schedule_id=%d", schedule_id)
	}

	if request_id > 0 {
		query += fmt.Sprintf(" AND w.request_id=%d", request_id)
	}
	log.Println(query)
	// prepare
	stmt, err := config.App().DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	// query
	rows, err := stmt.Query(user_id)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	webhooks := []*Webhook{}
	for rows.Next() {
		webhook := &Webhook{
			Schedule: &Schedule{},
			Request:  &Request{},
		}

		if err := rows.Scan(&webhook.ID, &webhook.ScheduleID, &webhook.RequestID, &webhook.Active, &webhook.CreatedAt, &webhook.UpdatedAt, &webhook.DeletedAt, &webhook.Schedule.Timing, &webhook.Request.Url); err != nil {
			return webhooks, err
		}

		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

func (m *Webhook) Paginate(userID, offset, limit int, search string) []*Webhook {
	webhooks := []*Webhook{}

	// prepare requests paginate
	stmt, err := config.App().DB.Prepare(config.App().QUERY["WEBHOOKS_PAGINATE"])
	if err != nil {
		return webhooks
	}

	// query
	rows, err := stmt.Query(userID, "%"+search+"%", offset, limit)
	if err != nil {
		return webhooks
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		webhook := &Webhook{
			Schedule: &Schedule{},
			Request:  &Request{},
		}

		if err := rows.Scan(&webhook.ID, &webhook.ScheduleID, &webhook.RequestID, &webhook.Active, &webhook.CreatedAt, &webhook.UpdatedAt, &webhook.DeletedAt, &webhook.Schedule.Timing, &webhook.Request.Url); err != nil {
			return webhooks
		}

		webhooks = append(webhooks, webhook)
	}

	return webhooks
}

func (m *Webhook) Create() error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["WEBHOOK_INSERT"])
	if err != nil {
		return err
	}

	// user_id,url,method,content,active
	err = stmt.QueryRow(m.ScheduleID, m.RequestID, m.Active).Scan(&m.ID, &m.ScheduleID, &m.RequestID, &m.Active)
	if err != nil {
		return err
	}

	return nil
}

func (m *Webhook) IDExists(id, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["WEBHOOK_ID_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(id, userID)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return false, err
		}
	}
	return exists > 0, nil
}

func (m *Webhook) UniqExists(scheduleID, requestID, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["WEBHOOK_UNIQ_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(scheduleID, requestID, userID)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return false, err
		}
	}
	return exists > 0, nil
}

func (m *Webhook) Update(query string, params []any) error {

	stmt, err := config.App().DB.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(params...)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("Webhook not updated")
	}

	return nil
}

func (m *Webhook) Delete(id, userID int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["WEBHOOK_DELETE"])
	if err != nil {
		return err
	}

	deleteAndUpdate := time.Now().Format("2006-01-02 15:04:05")

	result, err := stmt.Exec(deleteAndUpdate, deleteAndUpdate, id, userID)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("Webhook not deleted")
	}

	return nil
}
