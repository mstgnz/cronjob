package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/services"
)

type GroupHandler struct{}

func (h *GroupHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return services.Render(w, r, "group", map[string]any{}, "group/list", "group/new")
}
