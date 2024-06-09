package models

import (
	"time"

	"github.com/mstgnz/cronjob/config"
)

type AppLog struct {
	ID        int        `json:"id"`
	Level     string     `json:"level"`
	Message   string     `json:"message"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (m *AppLog) Create(level, message string) error {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["APP_LOG_INSERT"])
	if err != nil {
		return err
	}

	_, err = stmt.Exec(level, message)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	return nil
}
