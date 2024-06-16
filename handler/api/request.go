package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/services"
)

type RequestHandler struct {
	*services.RequestService
}

func (h *RequestHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.ListService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *RequestHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.CreateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *RequestHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.UpdateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *RequestHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.DeleteService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *RequestHandler) BulkHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.RequestBulkService(w, r)
	return config.WriteJSON(w, statusCode, response)
}
