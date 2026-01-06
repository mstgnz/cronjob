package web

import (
	"encoding/json"
	"net/http"

	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/pkg/config"
	"github.com/mstgnz/cronjob/pkg/load"
	"github.com/mstgnz/cronjob/pkg/response"
)

type HomeHandler struct{}

func (h *HomeHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	user, ok := r.Context().Value(config.CKey("user")).(*models.User)
	if !ok {
		return load.Render(w, r, "home", map[string]any{})
	}

	// Get statistics
	schedule := &models.Schedule{}
	group := &models.Group{}
	notification := &models.Notification{}
	request := &models.Request{}
	webhook := &models.Webhook{}
	scheduleLog := &models.ScheduleLog{}

	// Get daily counts for last 7 days
	dailyCounts := scheduleLog.DailyCounts(user.ID, 7)

	// Convert to JSON string for template
	dailyLogsJSON, _ := json.Marshal(dailyCounts)

	stats := map[string]any{
		"schedules":     schedule.Count(user.ID),
		"groups":        group.Count(user.ID),
		"notifications": notification.Count(user.ID),
		"requests":      request.Count(user.ID),
		"webhooks":      webhook.Count(user.ID),
		"logs":          scheduleLog.Count(user.ID),
		"dailyLogsJSON": string(dailyLogsJSON),
	}

	return load.Render(w, r, "home", stats)
}

func (h *HomeHandler) TriggerHandler(w http.ResponseWriter, r *http.Request) error {
	return response.WriteJSON(w, http.StatusOK, response.Response{Status: true, Message: "triggered", Data: nil})
}
