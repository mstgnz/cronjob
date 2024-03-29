package handler

import (
	"log"
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type Admin struct{}

func (a *Admin) HomeHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "home", map[string]any{}); err != nil {
		log.Println(err)
	}
}

func (a *Admin) ScheduleHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "home", map[string]any{}); err != nil {
		log.Println(err)
	}
}
