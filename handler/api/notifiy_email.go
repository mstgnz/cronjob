package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
)

type NotifyEmailHandler struct{}

func (h *NotifyEmailHandler) NotifyEmailListHandler(w http.ResponseWriter, r *http.Request) error {
	notifyEmail := &models.NotifyEmail{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	email := r.URL.Query().Get("email")

	notifyEmails, err := notifyEmail.Get(cUser.ID, id, email)
	if err != nil {
		return config.WriteJSON(w, http.StatusOK, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: notifyEmails})
}

func (h *NotifyEmailHandler) NotifyEmailCreateHandler(w http.ResponseWriter, r *http.Request) error {
	notifyEmail := &models.NotifyEmail{}
	if err := config.ReadJSON(w, r, notifyEmail); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(notifyEmail)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification := &models.Notification{}
	exists, err := notification.IDExists(notifyEmail.NotificationID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Notification not found"})
	}

	exists, err = notifyEmail.EmailExists(config.App().DB.DB, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if exists {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Email already exists"})
	}

	err = notifyEmail.Create(config.App().DB.DB)
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Notify email created", Data: notifyEmail})
}

func (h *NotifyEmailHandler) NotifyEmailUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	updateData := &models.NotifyEmailUpdate{}
	if err := config.ReadJSON(w, r, &updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(updateData)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notifyEmail := &models.NotifyEmail{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := notifyEmail.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Notify email not found"})
	}

	queryParts := []string{"UPDATE notify_emails SET"}
	params := []any{}
	paramCount := 1

	if updateData.NotificationID > 0 {
		queryParts = append(queryParts, fmt.Sprintf("notification_id=$%d,", paramCount))
		params = append(params, updateData.NotificationID)
		paramCount++
	}
	if updateData.Email != "" {
		queryParts = append(queryParts, fmt.Sprintf("email=$%d,", paramCount))
		params = append(params, updateData.Email)
		paramCount++
	}
	if updateData.Active != nil {
		queryParts = append(queryParts, fmt.Sprintf("active=$%d,", paramCount))
		params = append(params, updateData.Active)
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

	queryParts = append(queryParts, fmt.Sprintf("WHERE id=$%d", paramCount))
	params = append(params, id)
	query := strings.Join(queryParts, " ")

	err = notifyEmail.Update(query, params)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (h *NotifyEmailHandler) NotifyEmailDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notifyEmail := &models.NotifyEmail{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := notifyEmail.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Notify email not found"})
	}

	err = notifyEmail.Delete(id, cUser.ID)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Soft delte success"})
}
