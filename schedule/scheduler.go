package schedule

import (
	"log"
	"time"
	_ "time/tzdata"

	"github.com/mstgnz/cronjob/config"
	"github.com/robfig/cron/v3"
)

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
			log.Println("running schedulr every 1 minute")
		})

	}); err != nil {
		log.Println("AddFunc AlotechCall", err)
	}

}
