package schedule

import (
	"time"

	"github.com/mstgnz/cronjob/config"
)

func AlotechCall() {
	config.IncrementRunning()
	defer func() {
		config.DecrementRunning()
	}()
	config.App().InfoLog.Println("running AlotechCall")
	time.Sleep(time.Second * 70)
	config.App().InfoLog.Println("finish AlotechCall")
}
