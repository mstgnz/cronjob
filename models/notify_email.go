package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type NotifyEmail struct {
	ID             int        `json:"id"`
	NotificationID int        `json:"notification_id" validate:"required,number"`
	Email          string     `json:"email" validate:"required,email"`
	Active         bool       `json:"active" validate:"required,boolean"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type NotifyEmailUpdate struct {
	NotificationID int    `json:"notification_id" validate:"omitempty,number"`
	Email          string `json:"email" validate:"omitempty,email"`
	Active         *bool  `json:"active" validate:"omitnil,boolean"`
}

func (m *NotifyEmail) Get(userID, id int, email string) ([]NotifyEmail, error) {

	query := strings.TrimSuffix(config.App().QUERY["NOTIFICATION_EMAILS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND ne.id=%v", id)
	}
	if email != "" {
		query += fmt.Sprintf(" AND ne.email=%v", email)
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

	var notifyEmails []NotifyEmail
	for rows.Next() {
		var notifyEmail NotifyEmail
		if err := rows.Scan(&notifyEmail.ID, &notifyEmail.NotificationID, &notifyEmail.Email, &notifyEmail.Active, &notifyEmail.CreatedAt, &notifyEmail.UpdatedAt, &notifyEmail.DeletedAt); err != nil {
			return nil, err
		}
		notifyEmails = append(notifyEmails, notifyEmail)
	}

	return notifyEmails, nil
}

func (m *NotifyEmail) Create() error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_EMAIL_INSERT"])
	if err != nil {
		return err
	}

	// user_id,url,method,content,active
	err = stmt.QueryRow(m.NotificationID, m.Email, m.Active).Scan(&m.ID, &m.NotificationID, &m.Email, &m.Active)
	if err != nil {
		return err
	}

	return nil
}

func (m *NotifyEmail) EmailExists(userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_EMAIL_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(userID, m.Email, m.NotificationID)
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

func (m *NotifyEmail) IDExists(id, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_EMAIL_ID_EXISTS_WITH_USER"])
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

func (m *NotifyEmail) Update(query string, params []any) error {

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
		return fmt.Errorf("Notification email not updated")
	}

	return nil
}

func (m *NotifyEmail) Delete(id, userID int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_EMAIL_DELETE"])
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
		return fmt.Errorf("Notification email not deleted")
	}

	return nil
}
