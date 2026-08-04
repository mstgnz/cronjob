package schedule

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mstgnz/cronjob/pkg/config"
	"github.com/mstgnz/cronjob/pkg/logger"
	"github.com/robfig/cron/v3"
)

// defaultRetentionDays is how long log rows are kept when LOG_RETENTION_DAYS is unset.
const defaultRetentionDays = 90

// CallRetention registers the nightly cleanup. Both log tables grow with every
// execution and nothing else ever deletes from them.
func CallRetention(c *cron.Cron) {
	days := retentionDays()
	if days <= 0 {
		logger.Info("Log retention disabled")
		return
	}

	if _, err := c.AddFunc("0 3 * * *", func() {
		pruneLogs(days)
	}); err != nil {
		logger.Warn("Retention Schedule Error", err.Error())
	}
}

func retentionDays() int {
	value := os.Getenv("LOG_RETENTION_DAYS")
	if value == "" {
		return defaultRetentionDays
	}
	days, err := strconv.Atoi(value)
	if err != nil {
		logger.Warn("Retention Config Error", fmt.Sprintf("LOG_RETENTION_DAYS=%q is not a number, using %d", value, defaultRetentionDays))
		return defaultRetentionDays
	}
	return days
}

func pruneLogs(days int) {
	for _, query := range []string{"SCHEDULE_LOGS_PRUNE", "APP_LOGS_PRUNE"} {
		result, err := config.App().DB.Exec(config.App().QUERY[query], strconv.Itoa(days))
		if err != nil {
			logger.Warn("Retention Error", fmt.Sprintf("%s: %v", query, err))
			continue
		}
		if affected, err := result.RowsAffected(); err == nil && affected > 0 {
			logger.Info("Log retention", fmt.Sprintf("%s removed %d rows older than %d days", query, affected, days))
		}
	}
}
