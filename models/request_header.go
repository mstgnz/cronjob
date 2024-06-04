package models

import "time"

type RequestHeader struct {
	ID        int        `json:"id"`
	RequestID int        `json:"request_id"`
	Header    string     `json:"header"`
	Active    bool       `json:"active"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
