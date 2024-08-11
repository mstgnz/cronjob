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

type NotificationService struct{}

func (s *NotificationService) ListService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	notification := &models.Notification{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	title := r.URL.Query().Get("title")

	notifications, err := notification.Get(cUser.ID, id, title)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"notifications": notifications}}
}

func (s *NotificationService) CreateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	notification := &models.Notification{}
	if err := response.ReadJSON(w, r, notification); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(notification)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification.UserID = cUser.ID

	exists, err := notification.TitleExists()
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if exists {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Title already exists"}
	}

	err = notification.Create(config.App().DB.DB)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusCreated, response.Response{Status: true, Message: "Notification created", Data: map[string]any{"notification": notification}}
}

func (s *NotificationService) UpdateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	updateData := &models.NotificationUpdate{}
	if err := response.ReadJSON(w, r, &updateData); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(updateData)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification := &models.Notification{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := notification.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Notification not found"}
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
	if updateData.IsMessage != nil {
		queryParts = append(queryParts, fmt.Sprintf("is_message=$%d,", paramCount))
		params = append(params, updateData.IsMessage)
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

	queryParts = append(queryParts, fmt.Sprintf("WHERE id=$%d AND user_id=$%d", paramCount, paramCount+1))
	params = append(params, id, cUser.ID)
	query := strings.Join(queryParts, " ")

	err = notification.Update(query, params)

	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"update": updateData}}
}

func (s *NotificationService) DeleteService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification := &models.Notification{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := notification.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Notification not found"}
	}

	err = notification.Delete(id, cUser.ID)

	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Soft delte success"}
}

func (s *NotificationService) BulkService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	bulk := &models.NotificationBulk{}
	if err := response.ReadJSON(w, r, bulk); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(bulk)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification := &models.Notification{
		UserID:    cUser.ID,
		Title:     bulk.Title,
		Content:   bulk.Content,
		IsMail:    bulk.IsMail,
		IsMessage: bulk.IsMessage,
		Active:    bulk.Active,
	}

	exists, err := notification.TitleExists()
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if exists {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Title already exists"}
	}

	tx, err := config.App().DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	err = notification.Create(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
		}
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	for _, email := range bulk.NotifyEmails {
		notifyEmail := &models.NotifyEmail{
			NotificationID: notification.ID,
			Email:          email.Email,
			Active:         email.Active,
		}

		// check header key
		exists, err = notifyEmail.EmailExists(tx, cUser.ID)
		if err != nil || exists {
			continue
		}

		err = notifyEmail.Create(tx)
		if err != nil {
			continue
		}
	}

	for _, message := range bulk.NotifyMessages {
		notifyMessage := &models.NotifyMessage{
			NotificationID: notification.ID,
			Phone:          message.Phone,
			Active:         message.Active,
		}

		// check header key
		exists, err = notifyMessage.PhoneExists(tx, cUser.ID)
		if err != nil || exists {
			continue
		}

		err = notifyMessage.Create(tx)
		if err != nil {
			continue
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusCreated, response.Response{Status: true, Message: "Notification created", Data: map[string]any{"notification": notification}}
}
