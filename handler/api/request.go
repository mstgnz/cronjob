package api

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type RequestHandler struct{}

func (a *RequestHandler) RequestListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *RequestHandler) RequestCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *RequestHandler) RequestUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *RequestHandler) RequestDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *RequestHandler) RequestHeaderListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *RequestHandler) RequestHeaderCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *RequestHandler) RequestHeaderUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *RequestHandler) RequestHeaderDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}
