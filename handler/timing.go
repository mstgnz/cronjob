package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type Timing struct{}

func (s *Timing) ScheduleHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "schedule", map[string]any{"products": map[string]any{"test": template.HTML("<strong>test</strong>")}}, "navlink", "subscribe", "recommend", "scroll"); err != nil {
		log.Println(err)
	}
}
