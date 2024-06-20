package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type NotificationHandler struct{}

func (h *NotificationHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, "notification", map[string]any{})
}
