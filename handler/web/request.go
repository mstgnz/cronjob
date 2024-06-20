package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type RequestHandler struct{}

func (h *RequestHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, r, "request", map[string]any{}, "request/list", "request/header", "request/new")
}
