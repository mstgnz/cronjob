package models

import "time"

type Group struct {
	ID        int        `json:"id"`
	UID       int        `json:"uid"`
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (g *Group) GetLogs(offset, limit int) []*Group {
	Groups := []*Group{}
	return Groups
}

func (g *Group) GetLog(id int) []*Group {
	Groups := []*Group{}
	return Groups
}

func (g Group) CreateLog(id int, is_error bool, log string) Group {
	return g
}
