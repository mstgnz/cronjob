package models

import "time"

type ScheduleMail struct {
	ID         int        `json:"id"`
	ScheduleID int        `json:"schedule_id"`
	Email      string     `json:"email"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

func (sm ScheduleMail) GetScheduleMail(scheduleId int) ScheduleMail {
	return sm
}

func (sm ScheduleMail) UpdateScheduleMail(id, scheduleId int, email string) ScheduleMail {
	return sm
}

func (sm ScheduleMail) DeleteScheduleMail(id int) ScheduleMail {
	return sm
}

func (sm ScheduleMail) CreateScheduleMail(scheduleId int, email string) ScheduleMail {
	return sm
}
