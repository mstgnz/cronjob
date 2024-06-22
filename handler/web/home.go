package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/services"
)

type HomeHandler struct{}

func (h *HomeHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return services.Render(w, r, "home", map[string]any{})
}
