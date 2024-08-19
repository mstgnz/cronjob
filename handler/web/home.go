package web

import (
	"net/http"

	"github.com/mstgnz/cronjob/pkg/load"
	"github.com/mstgnz/cronjob/pkg/response"
)

type HomeHandler struct{}

func (h *HomeHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return load.Render(w, r, "home", map[string]any{})
}

func (h *HomeHandler) TriggerHandler(w http.ResponseWriter, r *http.Request) error {
	return response.WriteJSON(w, http.StatusOK, response.Response{Status: true, Message: "triggered", Data: nil})
}
