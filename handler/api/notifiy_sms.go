package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type NotifySmsHandler struct{}

func (h *NotifySmsHandler) NotifySmsListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (h *NotifySmsHandler) NotifySmsCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (h *NotifySmsHandler) NotifySmsUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (h *NotifySmsHandler) NotifySmsDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}
