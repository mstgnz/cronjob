package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/services"
)

type GroupHandler struct {
	*services.GroupService
}

func (h *GroupHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.ListService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *GroupHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.CreateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *GroupHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.UpdateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *GroupHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.DeleteService(w, r)
	return response.WriteJSON(w, statusCode, result)
}
