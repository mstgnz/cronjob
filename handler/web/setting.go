package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type SettingHandler struct{}

func (h *SettingHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, r, "setting", map[string]any{})
}
