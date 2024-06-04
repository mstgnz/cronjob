package models

import "time"

type NotifySms struct {
	ID             int        `json:"id"`
	NotificationID int        `json:"notification_id"`
	Phone          string     `json:"phone"`
	Active         bool       `json:"active"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
