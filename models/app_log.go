package models

import "time"

type AppLog struct {
	ID        int        `json:"id"`
	IsError   bool       `json:"is_error"`
	Log       string     `json:"log"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (al *AppLog) GetLogs(offset, limit int) []*AppLog {
	appLogs := []*AppLog{}
	return appLogs
}

func (al *AppLog) GetLog(id int) []*AppLog {
	appLogs := []*AppLog{}
	return appLogs
}

func (al AppLog) CreateLog(id int, is_error bool, log string) AppLog {
	return al
}
