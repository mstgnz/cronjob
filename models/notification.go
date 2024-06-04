package models

import "time"

type Notification struct {
	ID         int        `json:"id"`
	ScheduleID int        `json:"schedule_id"`
	IsSms      bool       `json:"is_sms"`
	IsMail     bool       `json:"is_mail"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Active     bool       `json:"active"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
