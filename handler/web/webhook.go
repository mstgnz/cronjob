package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type WebhookHandler struct{}

func (h *WebhookHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, "webhook", map[string]any{})
}
