package config

import (
	"log/slog"
	"strings"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) logToDB(level string, message string) {
	stmt, err := App().DB.Prepare(App().QUERY["APP_LOG_INSERT"])
	if err == nil {
		_, _ = stmt.Exec(level, message)
	}
	defer func() {
		_ = stmt.Close()
	}()
}

func (l *Logger) Info(message ...string) {
	msg := strings.Join(message, ", ")
	slog.Info(msg)
	l.logToDB("INFO", msg)
}

func (l *Logger) Error(message ...string) {
	msg := strings.Join(message, ", ")
	slog.Error(msg)
	l.logToDB("ERROR", msg)
}

func (l *Logger) Warn(message ...string) {
	msg := strings.Join(message, ", ")
	slog.Error(msg)
	l.logToDB("WARNING", msg)
}

func (l *Logger) Debug(message ...string) {
	msg := strings.Join(message, ", ")
	slog.Debug(msg)
	l.logToDB("DEBUG", msg)
}
