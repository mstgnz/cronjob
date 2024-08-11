package logger

import (
	"log/slog"
	"strings"

	"github.com/mstgnz/cronjob/pkg/config"
)

func logToDB(level string, message string) {
	stmt, err := config.App().DB.Prepare(config.App().QUERY["APP_LOG_INSERT"])
	if err == nil {
		_, _ = stmt.Exec(level, message)
	}
	defer func() {
		_ = stmt.Close()
	}()
}

func Info(message ...string) {
	msg := strings.Join(message, ", ")
	slog.Info(msg)
	logToDB("INFO", msg)
}

func Error(message ...string) {
	msg := strings.Join(message, ", ")
	slog.Error(msg)
	logToDB("ERROR", msg)
}

func Warn(message ...string) {
	msg := strings.Join(message, ", ")
	slog.Error(msg)
	logToDB("WARNING", msg)
}

func Debug(message ...string) {
	msg := strings.Join(message, ", ")
	slog.Debug(msg)
	logToDB("DEBUG", msg)
}
