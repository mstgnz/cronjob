package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/services"
)

type SettingHandler struct{}

func (h *SettingHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return services.Render(w, r, "setting", map[string]any{})
}
