package models

import "time"

type ScheduleLog struct {
	ID         int        `json:"id"`
	ScheduleID int        `json:"schedule_id"`
	Took       float32    `json:"took"`
	Result     any        `json:"result"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

func (sl *ScheduleLog) GetSchedules(offset, limit int) []*ScheduleLog {
	scheduleLogs := []*ScheduleLog{}
	return scheduleLogs
}

func (sl *ScheduleLog) GetSchedule(scheduleId int) []*ScheduleLog {
	scheduleLogs := []*ScheduleLog{}
	return scheduleLogs
}

func (sl ScheduleLog) CreateScheduleLog(scheduleId int, startedAt, finishedAt time.Time, took float32, result any) ScheduleLog {
	return sl
}
