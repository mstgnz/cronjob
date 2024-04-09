package schedule

import (
	"log"
	"time"

	"github.com/mstgnz/cronjob/config"
)

func AlotechCall() {
	config.IncrementRunning()
	defer func() {
		config.DecrementRunning()
	}()
	log.Println("running AlotechCall")
	time.Sleep(time.Second * 70)
	log.Println("finish AlotechCall")
}
