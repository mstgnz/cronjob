package models

import "time"

type Webhook struct {
	ID         int        `json:"id"`
	ScheduleID int        `json:"schedule_id"`
	RequestID  string     `json:"request_id"`
	Active     bool       `json:"active"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
