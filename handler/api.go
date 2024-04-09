package handler

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/service"
)

type Api struct {
	service.Api
}

func (a *Api) LoginHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.LoginService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.RegisterService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) UserHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.UserService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.UserUpdateService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) ScheduleListHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.ScheduleListService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) ScheduleCreateHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.ScheduleCreateService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) ScheduleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.ScheduleUpdateService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) ScheduleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.ScheduleDeleteService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) ScheduleMailListHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.ScheduleMailListService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) ScheduleMailCreateHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.ScheduleMailCreateService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) ScheduleMailUpdateHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.ScheduleMailUpdateService(w, r)
	_ = config.WriteJSON(w, status, response)
}

func (a *Api) ScheduleMailDeleteHandler(w http.ResponseWriter, r *http.Request) {
	status, response := a.Api.ScheduleMailDeleteService(w, r)
	_ = config.WriteJSON(w, status, response)
}
