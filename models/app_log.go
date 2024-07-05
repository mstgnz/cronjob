package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
)

type AppLog struct {
	ID        int        `json:"id"`
	Level     string     `json:"level"`
	Message   string     `json:"message"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (m *AppLog) Count() int {
	rowCount := 0

	// prepare count
	stmt, err := config.App().DB.Prepare(config.App().QUERY["APP_LOG_COUNT"])
	if err != nil {
		return rowCount
	}

	// query
	rows, err := stmt.Query()
	if err != nil {
		return rowCount
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()
	for rows.Next() {
		if err := rows.Scan(&rowCount); err != nil {
			return rowCount
		}
	}

	return rowCount
}

func (m *AppLog) Get(offset, limit int, level string) ([]*AppLog, error) {

	query := strings.TrimSuffix(config.App().QUERY["APP_LOG_PAGINATE"], ";")

	if level != "" {
		query += fmt.Sprintf(" AND level='%s'", level)
	}

	// prepare
	stmt, err := config.App().DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	// query
	rows, err := stmt.Query(offset, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	appLogs := []*AppLog{}
	for rows.Next() {
		appLog := &AppLog{}
		if err := rows.Scan(&appLog.ID, &appLog.Level, &appLog.Message, &appLog.CreatedAt); err != nil {
			return nil, err
		}
		appLogs = append(appLogs, appLog)
	}

	return appLogs, nil
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
