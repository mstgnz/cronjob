package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type NotificationHandler struct{}

func (a *NotificationHandler) NotificationListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotificationCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotificationUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotificationDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotifyEmailListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotifyEmailCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotifyEmailUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotifyEmailDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotifySmsListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotifySmsCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotifySmsUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *NotificationHandler) NotifySmsDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}
