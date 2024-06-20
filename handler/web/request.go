package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type RequestHandler struct{}

func (h *RequestHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, "request", map[string]any{})
}
