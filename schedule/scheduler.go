package schedule

import (
	"log"
	"time"
	_ "time/tzdata"

	"github.com/mstgnz/cronjob/config"
	"github.com/robfig/cron/v3"
)

// This application utilizes a database locking mechanism to prevent duplicate cronjob tasks
// when running multiple instances in a Kubernetes (k8s) environment.

// The instance that first adds a record to the schedule_logs table will be the one triggering
// the cronjob task. Subsequent instances won't be able to create a new record due to an existing one,
// resulting in an error and preventing them from executing the task.
// The microsecond differences between instances will automatically facilitate the locking mechanism.

// Nevertheless, it is strongly advised for users of this project to establish an additional control
// mechanism on their own systems.
func CallSchedule(c *cron.Cron) {
	// set location
	loc, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		log.Println(err)
	}

	cron.WithLocation(loc)

	// Alotech Call - every night at 23:59
	if _, err = c.AddFunc("* * * * *", func() {
		config.ShuttingWrapper(func() {
			AlotechCall()
		})

	}); err != nil {
		log.Println("AddFunc AlotechCall", err)
	}

}
