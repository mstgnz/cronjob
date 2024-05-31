package models

import "time"

type Schedule struct {
	ID        int        `json:"id"`
	UserId    int        `json:"user_id"`
	Timing    string     `json:"timing"`
	Active    bool       `json:"active"`
	Running   bool       `json:"running"`
	SendMail  bool       `json:"send_mail"`
	Url       string     `json:"url"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (s *Schedule) GetSchedules(offset, limit int) []*Schedule {
	schedules := []*Schedule{}
	return schedules
}

func (s *Schedule) GetScheduleWithUserId(userId int) *Schedule {
	return s
}

func (s *Schedule) CreateSchedule(id, user_id int, timing, path string, active, running, send_mail bool) *Schedule {
	return s
}

func (s *Schedule) UpdateSchedule(id, user_id int, timing, path string, active, running, send_mail bool) *Schedule {
	return s
}

func (s *Schedule) DeleteSchedule(id int) *Schedule {
	return s
}
