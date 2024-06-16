package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/services"
)

type RequestHeaderHandler struct {
	*services.RequestHeaderService
}

func (h *RequestHeaderHandler) ListHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.ListService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *RequestHeaderHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.CreateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *RequestHeaderHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.UpdateService(w, r)
	return config.WriteJSON(w, statusCode, response)
}

func (h *RequestHeaderHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.DeleteService(w, r)
	return config.WriteJSON(w, statusCode, response)
}
