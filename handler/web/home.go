package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
)

type HomeHandler struct{}

func (h *HomeHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)
	return config.Render(w, r, "home", map[string]any{"user": cUser})
}
