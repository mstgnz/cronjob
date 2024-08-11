package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/services"
)

type ScheduleHandler struct {
	*services.ScheduleService
}

func (h *ScheduleHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.ListService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *ScheduleHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.CreateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *ScheduleHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.UpdateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *ScheduleHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.DeleteService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *ScheduleHandler) LogListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.LogListService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *ScheduleHandler) BulkHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.BulkService(w, r)
	return response.WriteJSON(w, statusCode, result)
}
