package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/services"
)

type ScheduleHandler struct{}

func (h *ScheduleHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return services.Render(w, r, "schedule", map[string]any{}, "schedule/list", "schedule/log", "schedule/new")
}
