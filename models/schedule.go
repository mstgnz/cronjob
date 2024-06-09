package models

import "time"

type Schedule struct {
	ID        int        `json:"id"`
	UserId    int        `json:"user_id"`
	GroupID   int        `json:"group_id"`
	RequestID int        `json:"request_id"`
	Timing    string     `json:"timing"`
	Timeout   int        `json:"timeout"`
	Retries   int        `json:"retries"`
	Running   bool       `json:"running"`
	Active    bool       `json:"active"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ScheduleUpdate struct {
	GroupID   int    `json:"group_id" validate:"required"`
	RequestID int    `json:"request_id" validate:"required"`
	Timing    string `json:"timing" validate:"required"`
	Timeout   int    `json:"timeout"`
	Retries   int    `json:"retries"`
	Running   *bool  `json:"running"`
	Active    *bool  `json:"active"`
}

func (m *Schedule) Get(offset, limit int) []*Schedule {
	schedules := []*Schedule{}
	return schedules
}

func (m *Schedule) GetWithUserId(userId int) *Schedule {
	return m
}

func (m *Schedule) Create(id, user_id int, timing, path string, active, running, send_mail bool) *Schedule {
	return m
}

func (m *Schedule) Update(id, user_id int, timing, path string, active, running, send_mail bool) *Schedule {
	return m
}

func (m *Schedule) Delete(id int) *Schedule {
	return m
}
