package models

import "time"

type Schedule struct {
	ID        int       `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Url       string    `json:"url"`
	IsAdmin   bool      `json:"is_admin"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

func (s *Schedule) GetSchedules(offset, limit int) []*Schedule {
	schedules := []*Schedule{}
	return schedules
}

func (s *Schedule) GetScheduleWithUserId(userId int) *Schedule {
	return s
}

func (s *Schedule) UpdateSchedule(id, user_id int, timing, path string, active, running, send_mail bool) *Schedule {
	return s
}

func (s *Schedule) DeleteSchedule(id int) *Schedule {
	return s
}
