package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type Notification struct {
	ID          int            `json:"id"`
	UserID      int            `json:"user_id" validate:"number"`
	Title       string         `json:"title" validate:"required"`
	Content     string         `json:"content" validate:"required"`
	IsMail      bool           `json:"is_mail" validate:"boolean"`
	IsSms       bool           `json:"is_sms" validate:"boolean"`
	Active      bool           `json:"active" validate:"boolean"`
	NotifyEmail []*NotifyEmail `json:"emails,omitempty"`
	NotifySms   []*NotifySms   `json:"sms,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
}

type NotificationUpdate struct {
	UserID  int    `json:"user_id" validate:"omitempty,number"`
	Title   string `json:"title" validate:"omitempty"`
	Content string `json:"content" validate:"omitempty"`
	IsSms   *bool  `json:"is_sms" validate:"omitnil,boolean"`
	IsMail  *bool  `json:"is_mail" validate:"omitnil,boolean"`
	Active  *bool  `json:"active" validate:"omitnil,boolean"`
}

func (m *Notification) Get(userID, id int, title string) ([]Notification, error) {

	query := strings.TrimSuffix(config.App().QUERY["NOTIFICATIONS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND id=%v", id)
	}
	if title != "" {
		query += fmt.Sprintf(" AND title=%v", title)
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

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Title, &notification.Content, &notification.IsMail, &notification.IsSms, &notification.Active, &notification.CreatedAt, &notification.UpdatedAt, &notification.DeletedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (m *Notification) Create() error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["NOTIFICATION_INSERT"])
	if err != nil {
		return err
	}

	// user_id,url,method,content,active
	err = stmt.QueryRow(m.UserID, m.Title, m.Content, m.IsMail, m.IsSms, m.Active).Scan(&m.ID, &m.UserID, &m.Title, &m.Content, &m.IsMail, &m.IsSms, &m.Active)
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
