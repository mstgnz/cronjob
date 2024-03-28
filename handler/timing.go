package handler

import (
	"log"
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type Timing struct{}

func (s *Timing) ScheduleHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "schedule", map[string]any{}); err != nil {
		log.Println(err)
	}
}
