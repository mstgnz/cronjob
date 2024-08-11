package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/services"
)

type HomeHandler struct{}

func (h *HomeHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return services.Render(w, r, "home", map[string]any{})
}

func (h *HomeHandler) TriggerHandler(w http.ResponseWriter, r *http.Request) error {
	return response.WriteJSON(w, http.StatusOK, response.Response{Status: true, Message: "triggered", Data: nil})
}
