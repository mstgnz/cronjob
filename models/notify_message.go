package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type NotifyMessage struct {
	ID             int        `json:"id"`
	NotificationID int        `json:"notification_id" validate:"required,number"`
	Phone          string     `json:"phone" validate:"required,e164"`
	Active         bool       `json:"active" validate:"boolean"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type NotifyMessageUpdate struct {
	NotificationID int    `json:"notification_id" validate:"omitempty,number"`
	Phone          string `json:"phone" validate:"omitempty,e164"`
	Active         *bool  `json:"active" validate:"omitnil,boolean"`
}

type NotifyMessageBulk struct {
	NotificationID int    `json:"notification_id" validate:"number"`
	Phone          string `json:"phone" validate:"required,e164"`
	Active         bool   `json:"active" validate:"boolean"`
}

func (m *NotifyMessage) Get(userID, id int, phone string) ([]NotifyMessage, error) {

	query := strings.TrimSuffix(config.App().QUERY["NOTIFICATION_MESSAGES"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND ns.id=%v", id)
	}
	if phone != "" {
		query += fmt.Sprintf(" AND ns.phone=%v", phone)
	}

	// prepare
	stmt, err := config.App().DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	// query
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	var notifyMessages []NotifyMessage
	for rows.Next() {
		var notifyMessage NotifyMessage
		if err := rows.Scan(&notifyMessage.ID, &notifyMessage.NotificationID, &notifyMessage.Phone, &notifyMessage.Active, &notifyMessage.CreatedAt, &notifyMessage.UpdatedAt, &notifyMessage.DeletedAt); err != nil {
			return nil, err
		}
		notifyMessages = append(notifyMessages, notifyMessage)
	}

	return notifyMessages, nil
}

func (m *NotifyMessage) Create(exec any) error {
	stmt, err := config.App().DB.RunPrepare(exec, config.App().QUERY["NOTIFICATION_MESSAGE_INSERT"])
	if err != nil {
		return err
	}

	// user_id,url,method,content,active
	err = stmt.QueryRow(m.NotificationID, m.Phone, m.Active).Scan(&m.ID, &m.NotificationID, &m.Phone, &m.Active)
	if err != nil {
		return err
	}

	return nil
}

func (m *NotifyMessage) PhoneExists(exec any, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.RunPrepare(exec, config.App().QUERY["NOTIFICATION_MESSAGE_PHONE_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(userID, m.Phone, m.NotificationID)
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

func (m *NotifyMessage) IDExists(id, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_MESSAGE_ID_EXISTS_WITH_USER"])
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

func (m *NotifyMessage) Update(query string, params []any) error {

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
		return fmt.Errorf("Notification message not updated")
	}

	return nil
}

func (m *NotifyMessage) Delete(id, userID int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_MESSAGE_DELETE"])
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
		return fmt.Errorf("Notification message not deleted")
	}

	return nil
}
