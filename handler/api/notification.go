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
	"github.com/mstgnz/cronjob/services"
)

type NotificationHandler struct {
	*services.NotificationService
}

func (h *NotificationHandler) NotificationListHandler(w http.ResponseWriter, r *http.Request) error {
	notification := &models.Notification{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	title := r.URL.Query().Get("title")

	notifications, err := notification.Get(cUser.ID, id, title)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: notifications})
}

func (h *NotificationHandler) NotificationCreateHandler(w http.ResponseWriter, r *http.Request) error {
	notification := &models.Notification{}
	if err := config.ReadJSON(w, r, notification); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(notification)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification.UserID = cUser.ID

	exists, err := notification.TitleExists()
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if exists {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Title already exists"})
	}

	err = notification.Create(config.App().DB.DB)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Notification created", Data: notification})
}

func (h *NotificationHandler) NotificationUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	updateData := &models.NotificationUpdate{}
	if err := config.ReadJSON(w, r, &updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(updateData)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification := &models.Notification{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := notification.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Notification not found"})
	}

	queryParts := []string{"UPDATE notifications SET"}
	params := []any{}
	paramCount := 1

	if updateData.Title != "" {
		queryParts = append(queryParts, fmt.Sprintf("title=$%d,", paramCount))
		params = append(params, updateData.Title)
		paramCount++
	}
	if updateData.Content != "" {
		queryParts = append(queryParts, fmt.Sprintf("content=$%d,", paramCount))
		params = append(params, updateData.Content)
		paramCount++
	}
	if updateData.IsMail != nil {
		queryParts = append(queryParts, fmt.Sprintf("is_mail=$%d,", paramCount))
		params = append(params, updateData.IsMail)
		paramCount++
	}
	if updateData.IsSms != nil {
		queryParts = append(queryParts, fmt.Sprintf("is_sms=$%d,", paramCount))
		params = append(params, updateData.IsSms)
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

	queryParts = append(queryParts, fmt.Sprintf("WHERE id=$%d AND user_id=$%d", paramCount, paramCount+1))
	params = append(params, id, cUser.ID)
	query := strings.Join(queryParts, " ")

	err = notification.Update(query, params)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (h *NotificationHandler) NotificationDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification := &models.Notification{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := notification.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Notification not found"})
	}

	err = notification.Delete(id, cUser.ID)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Soft delte success"})
}

func (h *NotificationHandler) NotificationBulkHandler(w http.ResponseWriter, r *http.Request) error {
	statusCode, response := h.NotificationBulkService(w, r)
	return config.WriteJSON(w, statusCode, response)
}
