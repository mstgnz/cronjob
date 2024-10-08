package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/pkg/config"
)

type Notification struct {
	ID             int              `json:"id"`
	UserID         int              `json:"user_id" validate:"number"`
	Title          string           `json:"title" validate:"required"`
	Content        string           `json:"content" validate:"required"`
	IsMail         bool             `json:"is_mail" validate:"boolean"`
	IsMessage      bool             `json:"is_message" validate:"boolean"`
	Active         bool             `json:"active" validate:"boolean"`
	User           *User            `json:"user,omitempty"`
	NotifyEmails   []*NotifyEmail   `json:"emails,omitempty"`
	NotifyMessages []*NotifyMessage `json:"messages,omitempty"`
	CreatedAt      *time.Time       `json:"created_at,omitempty"`
	UpdatedAt      *time.Time       `json:"updated_at,omitempty"`
	DeletedAt      *time.Time       `json:"deleted_at,omitempty"`
}

type NotificationUpdate struct {
	UserID    int    `json:"user_id" validate:"omitempty,number"`
	Title     string `json:"title" validate:"omitempty"`
	Content   string `json:"content" validate:"omitempty"`
	IsMessage *bool  `json:"is_message" validate:"omitnil,boolean"`
	IsMail    *bool  `json:"is_mail" validate:"omitnil,boolean"`
	Active    *bool  `json:"active" validate:"omitnil,boolean"`
}

type NotificationBulk struct {
	UserID         int                  `json:"user_id" validate:"number"`
	Title          string               `json:"title" validate:"required"`
	Content        string               `json:"content" validate:"required"`
	IsMail         bool                 `json:"is_mail" validate:"boolean"`
	IsMessage      bool                 `json:"is_message" validate:"boolean"`
	Active         bool                 `json:"active" validate:"boolean"`
	NotifyEmails   []*NotifyEmailBulk   `json:"notify_emails" validate:"required_without=NotifyMessages,dive"`
	NotifyMessages []*NotifyMessageBulk `json:"notify_messages" validate:"required_without=NotifyEmails,dive"`
}

func (m *Notification) Count(userId int) int {
	rowCount := 0

	// prepare count
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATIONS_COUNT"])
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

func (m *Notification) Get(userID, id int, title string) ([]*Notification, error) {

	query := strings.TrimSuffix(config.App().QUERY["NOTIFICATIONS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND id=%d", id)
	}
	if title != "" {
		query += fmt.Sprintf(" AND title='%s'", title)
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

	notifications := []*Notification{}
	for rows.Next() {
		notification := &Notification{}
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Title, &notification.Content, &notification.IsMail, &notification.IsMessage, &notification.Active, &notification.CreatedAt, &notification.UpdatedAt, &notification.DeletedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (m *Notification) Paginate(userID, offset, limit int, search string) []*Notification {
	notifications := []*Notification{}

	// prepare notifications paginate
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATIONS_PAGINATE"])
	if err != nil {
		return notifications
	}

	// query
	rows, err := stmt.Query("%"+search+"%", userID, offset, limit)
	if err != nil {
		return notifications
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		notification := &Notification{
			User: &User{},
		}

		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Title, &notification.Content, &notification.IsMail, &notification.IsMessage, &notification.Active, &notification.CreatedAt, &notification.UpdatedAt, &notification.DeletedAt, &notification.User.Fullname); err != nil {
			return notifications
		}

		notifications = append(notifications, notification)
	}

	return notifications
}

func (m *Notification) Create(exec any) error {
	stmt, err := config.App().DB.RunPrepare(exec, config.App().QUERY["NOTIFICATION_INSERT"])
	if err != nil {
		return err
	}

	// user_id,url,method,content,active
	err = stmt.QueryRow(m.UserID, m.Title, m.Content, m.IsMail, m.IsMessage, m.Active).Scan(&m.ID, &m.UserID, &m.Title, &m.Content, &m.IsMail, &m.IsMessage, &m.Active)
	if err != nil {
		return err
	}

	return nil
}

func (m *Notification) TitleExists() (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_TITLE_EXISTS_WITH_USER"])
	if err != nil {
		return false, err
	}

	// query
	rows, err := stmt.Query(m.UserID, m.Title)
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

func (m *Notification) IDExists(id, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_ID_EXISTS_WITH_USER"])
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

func (m *Notification) Update(query string, params []any) error {

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
		return fmt.Errorf("Notification not updated")
	}

	return nil
}

func (m *Notification) Delete(id, userID int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_DELETE"])
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
		return fmt.Errorf("Notification not deleted")
	}

	return nil
}
