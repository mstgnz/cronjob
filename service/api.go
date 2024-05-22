package service

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
)

type Api struct{}

func (a *Api) LoginService(w http.ResponseWriter, r *http.Request) (int, config.Response) {

	login := &models.UserLogin{}
	if err := config.ReadJSON(w, r, login); err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Invalid Content"}
	}

	err := config.Validate(login)
	if err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Failed to process request", Data: err.Error()}
	}

	user := &models.User{}
	user = user.GetUserWithMail(login.Email)
	if user.Email == "" {
		return http.StatusUnauthorized, config.Response{Status: false, Message: "User Not Found"}
	}

	if !config.ComparePassword(user.Password, login.Password) {
		return http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Credentials"}
	}

	token, err := config.GenerateToken(user.ID)
	if err != nil {
		return http.StatusUnauthorized, config.Response{Status: false, Message: "Failed to process request"}
	}

	return http.StatusOK, config.Response{Status: true, Message: "Login successful", Data: token}
}

func (a *Api) RegisterService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	register := &models.UserRegister{}
	if err := config.ReadJSON(w, r, register); err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Invalid Credentials"}
	}

	err := config.Validate(register)
	if err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Failed to process request", Data: err.Error()}
	}

	user := &models.User{}
	user = user.GetUserWithMail(register.Email)
	if user.Email != "" {
		return http.StatusUnauthorized, config.Response{Status: false, Message: "Email Already Exists"}
	}

	user, err = user.CreateUser(register)
	if err != nil {
		return http.StatusCreated, config.Response{Status: false, Message: "Failed Create User"}
	}

	token, err := config.GenerateToken(user.ID)
	if err != nil {
		return http.StatusUnauthorized, config.Response{Status: false, Message: "Failed to process request"}
	}

	return http.StatusOK, config.Response{Status: true, Message: "Register successful", Data: token}
}

func (a *Api) UserService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success", Data: r.Context().Value("user")}
}

func (a *Api) UserUpdateService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success"}
}

func (a *Api) ScheduleListService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success"}
}

func (a *Api) ScheduleCreateService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success"}
}

func (a *Api) ScheduleUpdateService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success"}
}

func (a *Api) ScheduleDeleteService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success"}
}

func (a *Api) ScheduleMailListService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success"}
}

func (a *Api) ScheduleMailCreateService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success"}
}

func (a *Api) ScheduleMailUpdateService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success"}
}

func (a *Api) ScheduleMailDeleteService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	return http.StatusOK, config.Response{Status: true, Message: "Success"}
}
