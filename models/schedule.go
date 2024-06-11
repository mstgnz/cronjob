package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type Schedule struct {
	ID             int        `json:"id"`
	UserID         int        `json:"user_id" validate:"number"`
	GroupID        int        `json:"group_id" validate:"required,number"`
	RequestID      int        `json:"request_id" validate:"required,number"`
	NotificationID int        `json:"notification_id" validate:"required,number"`
	Timing         string     `json:"timing" validate:"required,cron"` // https://crontab.guru/
	Timeout        int        `json:"timeout" validate:"number"`
	Retries        int        `json:"retries" validate:"number"`
	Running        bool       `json:"running" validate:"boolean"`
	Active         bool       `json:"active" validate:"required,boolean"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type ScheduleUpdate struct {
	GroupID        int    `json:"group_id" validate:"omitempty,number"`
	RequestID      int    `json:"request_id" validate:"omitempty,number"`
	NotificationID int    `json:"notification_id" validate:"omitempty,number"`
	Timing         string `json:"timing" validate:"omitempty,cron"`
	Timeout        *int   `json:"timeout" validate:"omitnil,number"`
	Retries        *int   `json:"retries" validate:"omitnil,number"`
	Active         *bool  `json:"active" validate:"omitnil,boolean"`
}

func (m *Schedule) Get(id, userID, groupID, requestID, NotificationID int, timing string) ([]Schedule, error) {

	query := strings.TrimSuffix(config.App().QUERY["SCHEDULES"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND id=%v", id)
	}
	if userID > 0 {
		query += fmt.Sprintf(" AND user_id=%v", userID)
	}
	if groupID > 0 {
		query += fmt.Sprintf(" AND group_id=%v", userID)
	}
	if requestID > 0 {
		query += fmt.Sprintf(" AND request_id=%v", requestID)
	}
	if NotificationID > 0 {
		query += fmt.Sprintf(" AND notification_id=%v", NotificationID)
	}
	if timing != "" {
		query += fmt.Sprintf(" AND timing=%v", timing)
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

	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		if err := rows.Scan(&schedule.ID, &schedule.UserID, &schedule.GroupID, &schedule.RequestID, &schedule.NotificationID, &schedule.Timing, &schedule.Timeout, &schedule.Retries, &schedule.Running, &schedule.Active, &schedule.CreatedAt, &schedule.UpdatedAt, &schedule.DeletedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (m *Schedule) Create() error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["SCHEDULES_INSERT"])
	if err != nil {
		return err
	}

	// user_id,group_id,request_id,timing,timeout,retries,running,active;
	err = stmt.QueryRow(m.UserID, m.GroupID, m.RequestID, m.NotificationID, m.Timing, m.Timeout, m.Retries, m.Active).Scan(&m.ID, &m.UserID, &m.GroupID, &m.RequestID, &m.NotificationID, &m.Timing, &m.Timeout, &m.Retries, &m.Running, &m.Active)
	if err != nil {
		return err
	}

	return nil
}

func (m *Schedule) IDExists(id, userID int) (bool, error) {
	exists := 0

	// prepare
	stmt, err := config.App().DB.Prepare(config.App().QUERY["SCHEDULES_ID_EXISTS_WITH_USER"])
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

func (m *Schedule) Update(query string, params []any) error {

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
		return fmt.Errorf("Schedule not updated")
	}

	return nil
}

func (m *Schedule) Delete(id, userID int) error {

	stmt, err := config.App().DB.Prepare(config.App().QUERY["SCHEDULES_DELETE"])
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
		return fmt.Errorf("Schedule not deleted")
	}

	return nil
}
