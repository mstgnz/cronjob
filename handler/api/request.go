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
	req := &models.Request{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	url := r.URL.Query().Get("url")

	requests, err := req.Get(cUser.ID, id, url)
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

	err = request.Create()
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Request created", Data: request})
}

func (h *RequestHandler) RequestUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	var updateData map[string]any
	if err := config.ReadJSON(w, r, &updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request := &models.Request{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	groupExists, err := request.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !groupExists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
	}

	queryParts := []string{"UPDATE requests SET"}
	params := []any{}
	paramCount := 1

	value, exists := updateData["url"]
	if exists {
		queryParts = append(queryParts, fmt.Sprintf("url=$%d,", paramCount))
		params = append(params, value)
		paramCount++
	}
	value, exists = updateData["method"]
	if exists {
		queryParts = append(queryParts, fmt.Sprintf("method=$%d,", paramCount))
		params = append(params, value)
		paramCount++
	}
	value, exists = updateData["content"]
	if exists {
		queryParts = append(queryParts, fmt.Sprintf("content=$%d,", paramCount))
		params = append(params, value)
		paramCount++
	}
	value, exists = updateData["active"]
	if exists {
		queryParts = append(queryParts, fmt.Sprintf("active=$%d,", paramCount))
		params = append(params, value)
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
	groupExists, err := request.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !groupExists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
	}

	err = request.Delete(id, cUser.ID)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Soft delte success"})
}

func (h *RequestHandler) RequestHeaderListHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (h *RequestHandler) RequestHeaderCreateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (h *RequestHandler) RequestHeaderUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}

func (h *RequestHandler) RequestHeaderDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success"})
}
