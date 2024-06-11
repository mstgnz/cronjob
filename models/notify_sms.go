package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type NotifySms struct {
	ID             int        `json:"id"`
	NotificationID int        `json:"notification_id" validate:"required,number"`
	Phone          string     `json:"phone" validate:"required,e164"`
	Active         bool       `json:"active" validate:"required,boolean"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type NotifySmsUpdate struct {
	NotificationID int    `json:"notification_id" validate:"omitempty,number"`
	Phone          string `json:"phone" validate:"omitempty,e164"`
	Active         *bool  `json:"active" validate:"omitnil,boolean"`
}

func (m *NotifySms) Get(userID, id int, phone string) ([]NotifySms, error) {

	query := strings.TrimSuffix(config.App().QUERY["NOTIFICATION_SMS"], ";")

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

	var notifySms []NotifySms
	for rows.Next() {
		var notifyS NotifySms
		if err := rows.Scan(&notifyS.ID, &notifyS.NotificationID, &notifyS.Phone, &notifyS.Active, &notifyS.CreatedAt, &notifyS.UpdatedAt, &notifyS.DeletedAt); err != nil {
			return nil, err
		}
		notifySms = append(notifySms, notifyS)
	}

	return notifySms, nil
}

func (m *NotifySms) Create() error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_SMS_INSERT"])
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

func (m *NotifySms) PhoneExists(userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_SMS_PHONE_EXISTS_WITH_USER"])
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

func (m *NotifySms) IDExists(id, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_SMS_ID_EXISTS_WITH_USER"])
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

func (m *NotifySms) Update(query string, params []any) error {

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
		return fmt.Errorf("Notification sms not updated")
	}

	return nil
}

func (m *NotifySms) Delete(id, userID int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_SMS_DELETE"])
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
		return fmt.Errorf("Notification sms not deleted")
	}

	return nil
}
