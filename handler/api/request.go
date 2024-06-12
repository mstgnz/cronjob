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

type RequestHandler struct{}

func (h *RequestHandler) RequestListHandler(w http.ResponseWriter, r *http.Request) error {
	request := &models.Request{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	url := r.URL.Query().Get("url")

	requests, err := request.Get(cUser.ID, id, url)
	if err != nil {
		return config.WriteJSON(w, http.StatusOK, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: requests})
}

func (h *RequestHandler) RequestCreateHandler(w http.ResponseWriter, r *http.Request) error {
	request := &models.Request{}
	if err := config.ReadJSON(w, r, request); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(request)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request.UserID = cUser.ID

	exists, err := request.UrlExists()
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if exists {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Url already exists"})
	}

	err = request.Create(config.App().DB.DB)
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Request created", Data: request})
}

func (h *RequestHandler) RequestUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	updateData := &models.RequestUpdate{}
	if err := config.ReadJSON(w, r, &updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(updateData)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request := &models.Request{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := request.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
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

	err = request.Update(query, params)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (h *RequestHandler) RequestDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request := &models.Request{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := request.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
	}

	err = request.Delete(id, cUser.ID)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Soft delte success"})
}

func (h *RequestHandler) RequestBulkHandler(w http.ResponseWriter, r *http.Request) error {
	bulk := &models.RequestBulk{}
	if err := config.ReadJSON(w, r, bulk); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(bulk)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
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
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if exists {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Url already exists"})
	}

	tx, err := config.App().DB.Begin()
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	err = request.Create(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: err.Error()})
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
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "test", Data: bulk})
}
