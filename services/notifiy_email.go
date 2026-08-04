package services

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/pkg/config"
	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/pkg/validate"
)

type NotifyEmailService struct{}

func (s *NotifyEmailService) ListService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	notifyEmail := &models.NotifyEmail{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	email := r.URL.Query().Get("email")

	notifyEmails, err := notifyEmail.Get(cUser.ID, id, email)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"notify_emails": notifyEmails}}
}

func (s *NotifyEmailService) CreateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	notifyEmail := &models.NotifyEmail{}
	if err := response.ReadJSON(w, r, notifyEmail); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(notifyEmail)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification := &models.Notification{}
	exists, err := notification.IDExists(notifyEmail.NotificationID, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Notification not found"}
	}

	exists, err = notifyEmail.EmailExists(config.App().DB.DB, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if exists {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Email already exists"}
	}

	err = notifyEmail.Create(config.App().DB.DB)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusCreated, response.Response{Status: true, Message: "Notify email created", Data: map[string]any{"notify_email": notifyEmail}}
}

func (s *NotifyEmailService) UpdateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	updateData := &models.NotifyEmailUpdate{}
	if err := response.ReadJSON(w, r, &updateData); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(updateData)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notifyEmail := &models.NotifyEmail{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := notifyEmail.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Notify email not found"}
	}

	// Owning the notify email row is not enough: the new notification must belong to
	// the caller too, otherwise this row can be re-parented onto someone else's job
	// and start delivering that job's output to an address the attacker controls.
	if updateData.NotificationID > 0 {
		notification := &models.Notification{}
		exists, err = notification.IDExists(updateData.NotificationID, cUser.ID)
		if err != nil {
			return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
		}
		if !exists {
			return http.StatusNotFound, response.Response{Status: false, Message: "Notification not found"}
		}
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
		return http.StatusBadRequest, response.Response{Status: false, Message: "No fields to update"}
	}

	// update at
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	queryParts = append(queryParts, fmt.Sprintf("updated_at=$%d", paramCount))
	params = append(params, updatedAt)
	paramCount++

	// ownership is enforced in the statement itself, not only by the check above
	queryParts = append(queryParts, fmt.Sprintf("WHERE id=$%d AND notification_id IN (SELECT id FROM notifications WHERE user_id=$%d)", paramCount, paramCount+1))
	params = append(params, id, cUser.ID)
	query := strings.Join(queryParts, " ")

	err = notifyEmail.Update(query, params)

	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"update": updateData}}
}

func (s *NotifyEmailService) DeleteService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notifyEmail := &models.NotifyEmail{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := notifyEmail.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Notify email not found"}
	}

	err = notifyEmail.Delete(id, cUser.ID)

	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Soft delte success"}
}
