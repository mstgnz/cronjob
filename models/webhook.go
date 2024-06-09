package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type Webhook struct {
	ID         int        `json:"id"`
	ScheduleID int        `json:"schedule_id" validate:"required,number"`
	RequestID  int        `json:"request_id" validate:"required,number"`
	Active     bool       `json:"active" validate:"required,boolean"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type WebhookUpdate struct {
	ScheduleID int   `json:"schedule_id" validate:"omitempty,number"`
	RequestID  int   `json:"request_id" validate:"omitempty,number"`
	Active     *bool `json:"active" validate:"omitnil,boolean"`
}

func (m *Webhook) Get(id, schedule_id, request_id, user_id int) ([]Webhook, error) {

	query := strings.TrimSuffix(config.App().QUERY["WEBHOOKS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND id=%v", id)
	}

	if schedule_id > 0 {
		query += fmt.Sprintf(" AND schedule_id=%v", schedule_id)
	}

	if request_id > 0 {
		query += fmt.Sprintf(" AND request_id=%v", request_id)
	}

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

	var webhooks []Webhook
	for rows.Next() {
		var webhook Webhook
		if err := rows.Scan(&webhook.ID, &webhook.ScheduleID, &webhook.RequestID, &webhook.Active, &webhook.CreatedAt, &webhook.UpdatedAt, &webhook.DeletedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
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
