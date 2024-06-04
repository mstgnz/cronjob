package models

import "time"

type Request struct {
	ID        int        `json:"id"`
	Url       string     `json:"url"`
	Method    string     `json:"method"`
	Content   string     `json:"content"`
	Active    bool       `json:"active"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
