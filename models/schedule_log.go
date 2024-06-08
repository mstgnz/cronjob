package models

import "time"

type ScheduleLog struct {
	ID         int        `json:"id"`
	ScheduleID int        `json:"schedule_id" validate:"required"`
	Took       float32    `json:"took" validate:"required"`
	Result     any        `json:"result" validate:"required"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

func (m *ScheduleLog) GetSchedules(offset, limit int) []*ScheduleLog {
	scheduleLogs := []*ScheduleLog{}
	return scheduleLogs
}

func (m *ScheduleLog) GetSchedule(scheduleId int) []*ScheduleLog {
	scheduleLogs := []*ScheduleLog{}
	return scheduleLogs
}

func (m ScheduleLog) CreateScheduleLog(scheduleId int, startedAt, finishedAt time.Time, took float32, result any) ScheduleLog {
	return m
}
