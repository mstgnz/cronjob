package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type HomeHandler struct{}

func (h *HomeHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, r, "home", map[string]any{})
}
