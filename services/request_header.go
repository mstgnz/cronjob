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

type RequestHeaderService struct{}

func (s *RequestHeaderService) ListService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	requestHeader := &models.RequestHeader{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	requestID, _ := strconv.Atoi(r.URL.Query().Get("request_id"))
	key := r.URL.Query().Get("key")

	requestHeaders, err := requestHeader.Get(cUser.ID, id, requestID, key)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"request_headers": requestHeaders}}
}

func (s *RequestHeaderService) CreateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	requestHeader := &models.RequestHeader{}
	if err := response.ReadJSON(w, r, requestHeader); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(requestHeader)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	// check request
	request := &models.Request{}
	exists, err := request.IDExists(requestHeader.RequestID, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Request not found"}
	}

	// check header key
	exists, err = requestHeader.HeaderExists(config.App().DB.DB, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if exists {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Header already exists"}
	}

	err = requestHeader.Create(config.App().DB.DB)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusCreated, response.Response{Status: true, Message: "Request Header created", Data: map[string]any{"request_header": requestHeader}}
}

func (s *RequestHeaderService) UpdateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	updateData := &models.RequestHeaderUpdate{}
	if err := response.ReadJSON(w, r, &updateData); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: err.Error()}
	}

	err := validate.Validate(updateData)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	requestHeader := &models.RequestHeader{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := requestHeader.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Request Header not found"}
	}

	queryParts := []string{"UPDATE request_headers SET"}
	params := []any{}
	paramCount := 1

	if updateData.RequestID > 0 {
		queryParts = append(queryParts, fmt.Sprintf("request_id=$%d,", paramCount))
		params = append(params, updateData.RequestID)
		paramCount++
	}
	if updateData.Key != "" {
		queryParts = append(queryParts, fmt.Sprintf("key=$%d,", paramCount))
		params = append(params, updateData.Key)
		paramCount++
	}
	if updateData.Value != "" {
		queryParts = append(queryParts, fmt.Sprintf("value=$%d,", paramCount))
		params = append(params, updateData.Value)
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

	queryParts = append(queryParts, fmt.Sprintf("WHERE id=$%d", paramCount))
	params = append(params, id)
	query := strings.Join(queryParts, " ")

	err = requestHeader.Update(query, params)

	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"update": updateData}}
}

func (s *RequestHeaderService) DeleteService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	requestHeader := &models.RequestHeader{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := requestHeader.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, response.Response{Status: false, Message: "Request Header not found"}
	}

	err = requestHeader.Delete(id, cUser.ID)

	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, response.Response{Status: true, Message: "Soft delte success"}
}
