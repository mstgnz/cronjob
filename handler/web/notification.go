package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type NotificationHandler struct{}

func (h *NotificationHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, r, "notification", map[string]any{}, "notification/list", "notification/email", "notification/message", "notification/new")
}
