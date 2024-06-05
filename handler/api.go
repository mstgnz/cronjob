package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
)

type Api struct{}

func (a *Api) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	login := &models.UserLogin{}
	if err := config.ReadJSON(w, r, login); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Invalid credentials"})
	}

	err := config.Validate(login)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	user := &models.User{}
	err = user.GetWithMail(login.Email)
	if err != nil {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: err.Error()})
	}

	if !config.ComparePassword(user.Password, login.Password) {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid credentials"})
	}

	token, err := config.GenerateToken(user.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusOK, config.Response{Status: false, Message: "Failed to generate token"})
	}

	// update last_login
	user.UpdateLastLogin()

	data := make(map[string]any)
	data["token"] = token
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Login successful", Data: data})
}

func (a *Api) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	register := &models.UserRegister{}
	if err := config.ReadJSON(w, r, register); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Invalid credentials"})
	}

	err := config.Validate(register)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	user := &models.User{}
	userExists, err := user.Exists(register.Email)
	if err != nil {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: err.Error()})
	}
	if userExists {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Email already exists"})
	}

	err = user.Create(register)
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: err.Error()})
	}

	token, err := config.GenerateToken(user.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: "Failed to generate token"})
	}

	data := make(map[string]any)
	data["token"] = token
	data["user"] = user
	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "User created", Data: data})
}

func (a *Api) UserHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: r.Context().Value(config.CKey("user"))})
}

func (a *Api) UserUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	updateData := &models.User{}
	if err := config.ReadJSON(w, r, updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	// get auth user in context
	cUser, ok := r.Context().Value(config.CKey("user")).(*models.User)

	if !ok || cUser == nil || cUser.ID == 0 {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	updateData.ID = cUser.ID
	queryParts := []string{"UPDATE users SET"}
	params := []any{}
	paramCount := 1

	if updateData.Fullname != "" {
		queryParts = append(queryParts, fmt.Sprintf("fullname=$%d,", paramCount))
		params = append(params, updateData.Fullname)
		paramCount++
	}
	if updateData.Email != "" {
		queryParts = append(queryParts, fmt.Sprintf("email=$%d,", paramCount))
		params = append(params, updateData.Email)
		paramCount++
	}
	if updateData.Phone != "" {
		queryParts = append(queryParts, fmt.Sprintf("phone=$%d,", paramCount))
		params = append(params, updateData.Phone)
		paramCount++
	}

	if len(params) == 0 {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "No fields to update"})
	}

	// update at
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	queryParts = append(queryParts, fmt.Sprintf("updated_at=$%d", paramCount))
	params = append(params, updatedAt)
	paramCount++

	// remove last comma
	// size := len(queryParts) - 1
	// queryParts[size] = strings.TrimSuffix(queryParts[size], ",")

	queryParts = append(queryParts, fmt.Sprintf("WHERE id=$%d", paramCount))
	params = append(params, updateData.ID)
	query := strings.Join(queryParts, " ")

	err := updateData.Update(query, params)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (a *Api) UserPassUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	updateData := &models.UserPasswordUpdate{}
	if err := config.ReadJSON(w, r, updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	err := config.Validate(updateData)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	if updateData.Password != updateData.RePassword {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Passwords do not match"})
	}

	// get auth user in context
	cUser, ok := r.Context().Value(config.CKey("user")).(*models.User)

	if !ok || cUser == nil || cUser.ID == 0 {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	user := &models.User{}
	user.ID = cUser.ID

	err = user.UpdatePassword(updateData.Password)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (a *Api) GroupListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) GroupCreateHandler(w http.ResponseWriter, r *http.Request) error {
	group := &models.Group{}
	if err := config.ReadJSON(w, r, group); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Invalid contents"})
	}

	err := config.Validate(group)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, ok := r.Context().Value(config.CKey("user")).(*models.User)

	if !ok || cUser == nil || cUser.ID == 0 {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	group.UserID = cUser.ID

	groupExists, err := group.Exists()
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if groupExists {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Group already exists"})
	}

	lastInsertId, err := group.Create()
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Group created", Data: lastInsertId})
}

func (a *Api) GroupUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) GroupDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) RequestListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) RequestCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) RequestUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) RequestDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) RequestHeaderListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) RequestHeaderCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) RequestHeaderUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) RequestHeaderDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotificationListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotificationCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotificationUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotificationDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotifyEmailListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotifyEmailCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotifyEmailUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotifyEmailDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotifySmsListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotifySmsCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotifySmsUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) NotifySmsDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) WebhookListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) WebhookCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) WebhookUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) WebhookDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
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

func (a *Api) ScheduleLogListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleLogCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleLogUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (a *Api) ScheduleLogDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}
