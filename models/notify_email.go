package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/pkg/config"
)

type NotifyEmail struct {
	ID             int           `json:"id"`
	NotificationID int           `json:"notification_id" validate:"required,number"`
	Email          string        `json:"email" validate:"required,email"`
	Active         bool          `json:"active" validate:"boolean"`
	Notification   *Notification `json:"notification,omitempty"`
	CreatedAt      *time.Time    `json:"created_at,omitempty"`
	UpdatedAt      *time.Time    `json:"updated_at,omitempty"`
	DeletedAt      *time.Time    `json:"deleted_at,omitempty"`
}

type NotifyEmailUpdate struct {
	NotificationID int    `json:"notification_id" validate:"omitempty,number"`
	Email          string `json:"email" validate:"omitempty,email"`
	Active         *bool  `json:"active" validate:"omitnil,boolean"`
}

type NotifyEmailBulk struct {
	NotificationID int    `json:"notification_id" validate:"number"`
	Email          string `json:"email" validate:"required,email"`
	Active         bool   `json:"active" validate:"boolean"`
}

func (m *NotifyEmail) Count(userId int) int {
	rowCount := 0

	// prepare count
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_EMAILS_COUNT"])
	if err != nil {
		return rowCount
	}

	// query
	rows, err := stmt.Query(userId)
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

func (m *NotifyEmail) Get(userID, id int, email string) ([]*NotifyEmail, error) {

	query := strings.TrimSuffix(config.App().QUERY["NOTIFICATION_EMAILS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND ne.id=%d", id)
	}
	if email != "" {
		query += fmt.Sprintf(" AND ne.email='%s'", email)
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

	notifyEmails := []*NotifyEmail{}
	for rows.Next() {
		notifyEmail := &NotifyEmail{
			Notification: &Notification{},
		}
		if err := rows.Scan(&notifyEmail.ID, &notifyEmail.NotificationID, &notifyEmail.Email, &notifyEmail.Active, &notifyEmail.CreatedAt, &notifyEmail.UpdatedAt, &notifyEmail.DeletedAt, &notifyEmail.Notification.Title); err != nil {
			return nil, err
		}
		notifyEmails = append(notifyEmails, notifyEmail)
	}

	return notifyEmails, nil
}

func (m *NotifyEmail) Paginate(userID, offset, limit int, search string) []*NotifyEmail {
	notifyEmails := []*NotifyEmail{}

	// prepare notify_emails paginate
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_EMAILS_PAGINATE"])
	if err != nil {
		return notifyEmails
	}

	// query
	rows, err := stmt.Query("%"+search+"%", userID, offset, limit)
	if err != nil {
		return notifyEmails
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		notifyEmail := &NotifyEmail{
			Notification: &Notification{},
		}

		if err := rows.Scan(&notifyEmail.ID, &notifyEmail.NotificationID, &notifyEmail.Email, &notifyEmail.Active, &notifyEmail.CreatedAt, &notifyEmail.UpdatedAt, &notifyEmail.DeletedAt, &notifyEmail.Notification.Title); err != nil {
			return notifyEmails
		}

		notifyEmails = append(notifyEmails, notifyEmail)
	}

	return notifyEmails
}

func (m *NotifyEmail) Create(exec any) error {
	stmt, err := config.App().DB.RunPrepare(exec, config.App().QUERY["NOTIFICATION_EMAIL_INSERT"])
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

func (m *NotifyEmail) EmailExists(exec any, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.RunPrepare(exec, config.App().QUERY["NOTIFICATION_EMAIL_EXISTS_WITH_USER"])
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
