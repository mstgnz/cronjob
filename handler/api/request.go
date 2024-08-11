package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/services"
)

type RequestHandler struct {
	*services.RequestService
}

func (h *RequestHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.ListService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *RequestHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.CreateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *RequestHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.UpdateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *RequestHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.DeleteService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *RequestHandler) BulkHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.RequestBulkService(w, r)
	return response.WriteJSON(w, statusCode, result)
}
