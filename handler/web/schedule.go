package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type ScheduleHandler struct{}

func (h *ScheduleHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, "schedule", map[string]any{})
}
