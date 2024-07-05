package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type ScheduleLog struct {
	ID         int        `json:"id"`
	ScheduleID int        `json:"schedule_id" validate:"required"`
	Took       float32    `json:"took" validate:"required"`
	Result     any        `json:"result" validate:"required"`
	Schedule   *Schedule  `json:"schedule,omitempty"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}

func (m *ScheduleLog) Count(userID int) int {
	rowCount := 0

	// prepare count
	stmt, err := config.App().DB.Prepare(config.App().QUERY["SCHEDULE_LOGS_COUNT"])
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

func (m *ScheduleLog) Get(id, schedule_id, user_id int) ([]*ScheduleLog, error) {
	query := strings.TrimSuffix(config.App().QUERY["SCHEDULE_LOGS"], ";")

	if id > 0 {
		query += fmt.Sprintf(" AND id=%v", id)
	}

	// prepare
	stmt, err := config.App().DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	// query
	rows, err := stmt.Query(schedule_id, user_id)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	scheduleLogs := []*ScheduleLog{}
	for rows.Next() {
		scheduleLog := &ScheduleLog{}
		if err := rows.Scan(&scheduleLog.ID, &scheduleLog.ScheduleID, &scheduleLog.StartedAt, &scheduleLog.FinishedAt, &scheduleLog.Took, &scheduleLog.Result, &scheduleLog.CreatedAt); err != nil {
			return nil, err
		}
		scheduleLogs = append(scheduleLogs, scheduleLog)
	}

	return scheduleLogs, nil
}

func (m *ScheduleLog) Paginate(userID, offset, limit int, search string) []*ScheduleLog {
	scheduleLogs := []*ScheduleLog{}

	// prepare paginate
	stmt, err := config.App().DB.Prepare(config.App().QUERY["SCHEDULE_LOGS_PAGINATE"])
	if err != nil {
		return scheduleLogs
	}

	// query
	rows, err := stmt.Query(userID, "%"+search+"%", offset, limit)
	if err != nil {
		return scheduleLogs
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		scheduleLog := &ScheduleLog{
			Schedule: &Schedule{},
		}

		if err := rows.Scan(&scheduleLog.ID, &scheduleLog.ScheduleID, &scheduleLog.StartedAt, &scheduleLog.FinishedAt, &scheduleLog.Took, &scheduleLog.Result, &scheduleLog.CreatedAt, &scheduleLog.Schedule.Timing); err != nil {
			return scheduleLogs
		}

		scheduleLogs = append(scheduleLogs, scheduleLog)
	}

	return scheduleLogs
}

func (m *ScheduleLog) Create(scheduleId int) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["SCHEDULE_LOG_INSERT"])
	if err != nil {
		return err
	}

	// user_id,url,method,content,active
	err = stmt.QueryRow(m.ScheduleID, m.StartedAt, m.FinishedAt, m.Took, m.Result).Scan(&m.ID, &m.ScheduleID, &m.StartedAt, &m.FinishedAt, &m.Took, &m.Result, &m.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
