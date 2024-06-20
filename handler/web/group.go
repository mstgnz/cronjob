package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type GroupHandler struct{}

func (h *GroupHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return config.Render(w, "group", map[string]any{}, "group/list", "group/new")
}
