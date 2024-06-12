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

type RequestHeaderHandler struct{}

func (h *RequestHeaderHandler) RequestHeaderListHandler(w http.ResponseWriter, r *http.Request) error {
	requestHeader := &models.RequestHeader{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	requestID, _ := strconv.Atoi(r.URL.Query().Get("request_id"))
	key := r.URL.Query().Get("key")

	requests, err := requestHeader.Get(cUser.ID, id, requestID, key)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: requests})
}

func (h *RequestHeaderHandler) RequestHeaderCreateHandler(w http.ResponseWriter, r *http.Request) error {
	requestHeader := &models.RequestHeader{}
	if err := config.ReadJSON(w, r, requestHeader); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(requestHeader)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	// check request
	request := &models.Request{}
	exists, err := request.IDExists(requestHeader.RequestID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
	}

	// check header key
	exists, err = requestHeader.HeaderExists(config.App().DB.DB, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if exists {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Header already exists"})
	}

	err = requestHeader.Create(config.App().DB.DB)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Request Header created", Data: requestHeader})
}

func (h *RequestHeaderHandler) RequestHeaderUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	updateData := &models.RequestHeaderUpdate{}
	if err := config.ReadJSON(w, r, &updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(updateData)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	requestHeader := &models.RequestHeader{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := requestHeader.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request Header not found"})
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

	err = requestHeader.Update(query, params)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (h *RequestHeaderHandler) RequestHeaderDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	requestHeader := &models.RequestHeader{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := requestHeader.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request Header not found"})
	}

	err = requestHeader.Delete(id, cUser.ID)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Soft delte success"})
}
