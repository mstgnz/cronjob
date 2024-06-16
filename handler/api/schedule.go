package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/services"
)

type ScheduleHandler struct {
	*services.ScheduleService
}

func (h *ScheduleHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.ListService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *ScheduleHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.CreateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *ScheduleHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.UpdateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *ScheduleHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.DeleteService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *ScheduleHandler) LogListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.LogListService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *ScheduleHandler) BulkHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.BulkService(w, r)
	return config.WriteJSON(w, statusCode, response)
}
