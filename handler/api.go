package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
)

type Api struct{}

func (a *Api) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	login := &models.UserLogin{}
	if err := config.ReadJSON(w, r, login); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Invalid Content"})
	}

	err := config.Validate(login)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Failed to process request", Data: err.Error()})
	}

	user := &models.User{}
	user = user.GetUserWithMail(login.Email)
	if user.Email == "" {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "User Not Found"})
	}

	if !config.ComparePassword(user.Password, login.Password) {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	token, err := config.GenerateToken(user.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Failed to process request"})
	}

	// update last_login
	user.UpdateLastLogin()

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Login successful", Data: token})
}

func (a *Api) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	register := &models.UserRegister{}
	if err := config.ReadJSON(w, r, register); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	err := config.Validate(register)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Failed to process request", Data: err.Error()})
	}

	user := &models.User{}
	user = user.GetUserWithMail(register.Email)
	if user.Email != "" {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Email Already Exists"})
	}

	user, err = user.CreateUser(register)
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: "Failed Create User"})
	}

	token, err := config.GenerateToken(user.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Failed to process request"})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Register successful", Data: token})
}

func (a *Api) UserHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: r.Context().Value(config.CKey("user"))})
}

func (a *Api) UserUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	updateData := &models.User{}
	if err := config.ReadJSON(w, r, updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	queryParts := []string{"UPDATE users SET"}
	params := []any{}
	paramCount := 1

	if updateData.Fullname != "" {
		queryParts = append(queryParts, fmt.Sprintf("fullname=$%d", paramCount))
		params = append(params, updateData.Fullname)
		paramCount++
	}
	if updateData.Email != "" {
		queryParts = append(queryParts, fmt.Sprintf("email=$%d", paramCount))
		params = append(params, updateData.Email)
		paramCount++
	}

	if len(params) == 0 {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "No fields to update"})
	}

	queryParts = append(queryParts, fmt.Sprintf("WHERE id=$%d", paramCount))
	params = append(params, updateData.ID)
	query := strings.Join(queryParts, " ")

	updateData = updateData.Update(query, params)

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (a *Api) ScheduleListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleMailListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleMailCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleMailUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleMailDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}
