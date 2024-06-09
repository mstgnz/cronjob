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

type WebhookHandler struct{}

func (h *WebhookHandler) WebhookListHandler(w http.ResponseWriter, r *http.Request) error {
	webhook := &models.Webhook{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	schedule_id, _ := strconv.Atoi(r.URL.Query().Get("schedule_id"))
	request_id, _ := strconv.Atoi(r.URL.Query().Get("request_id"))

	if schedule_id == 0 {
		return config.WriteJSON(w, http.StatusOK, config.Response{Status: false, Message: "schedule_id param required"})
	}

	requests, err := webhook.Get(id, schedule_id, request_id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusOK, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: requests})
}

func (h *WebhookHandler) WebhookCreateHandler(w http.ResponseWriter, r *http.Request) error {
	webhook := &models.Webhook{}
	if err := config.ReadJSON(w, r, webhook); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(webhook)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	// check schedule_id
	schedule := &models.Schedule{}
	exists, err := schedule.IDExists(webhook.ScheduleID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Schedule not found"})
	}

	// check request_id
	request := &models.Request{}
	exists, err = request.IDExists(webhook.RequestID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
	}

	// check schedule_id and request_id
	exists, err = webhook.UniqExists(webhook.ScheduleID, webhook.RequestID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Webhook already exists"})
	}

	err = webhook.Create()
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Webhook created", Data: webhook})
}

func (h *WebhookHandler) WebhookUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	updateData := &models.WebhookUpdate{}
	if err := config.ReadJSON(w, r, &updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(updateData)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	webhook := &models.Webhook{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := webhook.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Webhook not found"})
	}

	queryParts := []string{"UPDATE webhooks SET"}
	params := []any{}
	paramCount := 1

	if updateData.ScheduleID > 0 {
		queryParts = append(queryParts, fmt.Sprintf("schedule_id=$%d,", paramCount))
		params = append(params, updateData.ScheduleID)
		paramCount++
	}
	if updateData.RequestID > 0 {
		queryParts = append(queryParts, fmt.Sprintf("request_id=$%d,", paramCount))
		params = append(params, updateData.RequestID)
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

	err = webhook.Update(query, params)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (h *WebhookHandler) WebhookDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	webhook := &models.Webhook{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := webhook.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Webhook not found"})
	}

	err = webhook.Delete(id, cUser.ID)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Soft delte success"})
}
