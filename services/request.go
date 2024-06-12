package services

import (
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
)

type RequestService struct{}

func (s *RequestService) RequestBulkService(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	bulk := &models.RequestBulk{}
	if err := config.ReadJSON(w, r, bulk); err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: err.Error()}
	}

	err := config.Validate(bulk)
	if err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request := &models.Request{
		UserID:  cUser.ID,
		Url:     bulk.Url,
		Method:  bulk.Method,
		Content: bulk.Content,
		Active:  bulk.Active,
	}

	exists, err := request.UrlExists()
	if err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}
	if exists {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Url already exists"}
	}

	tx, err := config.App().DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}

	err = request.Create(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
		}
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}

	for _, header := range bulk.RequestHeaders {
		requestHeader := &models.RequestHeader{
			RequestID: request.ID,
			Key:       header.Key,
			Value:     header.Value,
			Active:    header.Active,
		}

		// check header key
		exists, err = requestHeader.HeaderExists(tx, cUser.ID)
		if err != nil || exists {
			continue
		}

		err = requestHeader.Create(tx)
		if err != nil {
			continue
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}

	return http.StatusCreated, config.Response{Status: true, Message: "test", Data: bulk}
}
