package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type NotifyEmailHandler struct{}

func (h *NotifyEmailHandler) NotifyEmailListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (h *NotifyEmailHandler) NotifyEmailCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (h *NotifyEmailHandler) NotifyEmailUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (h *NotifyEmailHandler) NotifyEmailDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}
