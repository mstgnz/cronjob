package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type WebhookHandler struct{}

func (a *WebhookHandler) WebhookListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *WebhookHandler) WebhookCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *WebhookHandler) WebhookUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *WebhookHandler) WebhookDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}
