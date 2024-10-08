package services

import (
	"encoding/json"
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

type RequestService struct{}

func (h *RequestService) ListService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	request := &models.Request{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	url := r.URL.Query().Get("url")

	requests, err := request.Get(cUser.ID, id, url)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"requests": requests}}
}

func (h *RequestService) CreateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	request := &models.Request{}
	if err := response.ReadJSON(w, r, request); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(request)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request.UserID = cUser.ID

	exists, err := request.UrlExists()
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if exists {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Url already exists"}
	}

	err = request.Create(config.App().DB.DB)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusCreated, response.Response{Status: true, Message: "Request created", Data: map[string]any{"request": request}}
}

func (h *RequestService) UpdateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	updateData := &models.RequestUpdate{}
	if err := response.ReadJSON(w, r, &updateData); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(updateData)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request := &models.Request{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := request.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Request not found"}
	}

	queryParts := []string{"UPDATE requests SET"}
	params := []any{}
	paramCount := 1

	if updateData.Url != "" {
		queryParts = append(queryParts, fmt.Sprintf("url=$%d,", paramCount))
		params = append(params, updateData.Url)
		paramCount++
	}
	if updateData.Method != "" {
		queryParts = append(queryParts, fmt.Sprintf("method=$%d,", paramCount))
		params = append(params, updateData.Method)
		paramCount++
	}
	if updateData.Content != "" {
		queryParts = append(queryParts, fmt.Sprintf("content=$%d,", paramCount))
		params = append(params, updateData.Content)
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

	err = request.Update(query, params)

	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"update": updateData}}
}

func (h *RequestService) DeleteService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request := &models.Request{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := request.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Request not found"}
	}

	err = request.Delete(id, cUser.ID)

	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Soft delte success"}
}

func (s *RequestService) RequestBulkService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	bulk := &models.RequestBulk{}
	if err := response.ReadJSON(w, r, bulk); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(bulk)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request := &models.Request{
		UserID:  cUser.ID,
		Url:     bulk.Url,
		Method:  bulk.Method,
		Content: json.RawMessage(bulk.Content),
		Active:  bulk.Active,
	}

	exists, err := request.UrlExists()
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if exists {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Url already exists"}
	}

	tx, err := config.App().DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	err = request.Create(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
		}
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
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
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusCreated, response.Response{Status: true, Message: "test", Data: map[string]any{"request": request}}
}
