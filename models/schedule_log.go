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
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}

func (m *ScheduleLog) GetSchedules(id, schedule_id, user_id int) ([]ScheduleLog, error) {
	query := strings.TrimSuffix(config.App().QUERY["REQUESTS"], ";")

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

	var scheduleLogs []ScheduleLog
	for rows.Next() {
		var scheduleLog ScheduleLog
		if err := rows.Scan(&scheduleLog.ID, &scheduleLog.ScheduleID, &scheduleLog.StartedAt, &scheduleLog.FinishedAt, &scheduleLog.Took, &scheduleLog.Result, &scheduleLog.CreatedAt); err != nil {
			return nil, err
		}
		scheduleLogs = append(scheduleLogs, scheduleLog)
	}

	return scheduleLogs, nil
}

func (m *ScheduleLog) Create(scheduleId int) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["REQUEST_INSERT"])
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
