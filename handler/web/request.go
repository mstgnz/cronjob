package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/services"
)

type RequestHandler struct{}

func (h *RequestHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return services.Render(w, r, "request", map[string]any{}, "request/list", "request/header", "request/new")
}
