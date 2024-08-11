package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/services"
)

type RequestHeaderHandler struct {
	*services.RequestHeaderService
}

func (h *RequestHeaderHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.ListService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *RequestHeaderHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.CreateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *RequestHeaderHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.UpdateService(w, r)
	return response.WriteJSON(w, statusCode, result)
}

func (h *RequestHeaderHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, result := h.DeleteService(w, r)
	return response.WriteJSON(w, statusCode, result)
}
