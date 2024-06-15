package services

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
)

type NotificationService struct{}

func (s *NotificationService) NotificationBulkService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	bulk := &models.NotificationBulk{}
	if err := config.ReadJSON(w, r, bulk); err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: err.Error()}
	}

	err := config.Validate(bulk)
	if err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()}
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
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}
	if exists {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Title already exists"}
	}

	tx, err := config.App().DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}

	err = notification.Create(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
		}
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
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
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}

	return http.StatusCreated, config.Response{Status: true, Message: "Notification created", Data: notification}
}
