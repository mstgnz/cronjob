package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/services"
)

type GroupHandler struct {
	*services.GroupService
}

func (h *GroupHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.ListSerice(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *GroupHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.CreateSerice(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *GroupHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.UpdateSerice(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *GroupHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.DeleteSerice(w, r)
	return config.WriteJSON(w, statusCode, response)
}
