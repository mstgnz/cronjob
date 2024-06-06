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

type GroupHandler struct{}

func (h *GroupHandler) GroupListHandler(w http.ResponseWriter, r *http.Request) error {
	group := &models.Group{}

	// get auth user in context
	cUser, ok := r.Context().Value(config.CKey("user")).(*models.User)

	if !ok || cUser == nil || cUser.ID == 0 {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))

	groups, err := group.Get(cUser.ID, id, uid)
	if err != nil {
		return config.WriteJSON(w, http.StatusOK, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: groups})
}

func (h *GroupHandler) GroupCreateHandler(w http.ResponseWriter, r *http.Request) error {
	group := &models.Group{}
	if err := config.ReadJSON(w, r, group); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(group)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, ok := r.Context().Value(config.CKey("user")).(*models.User)

	if !ok || cUser == nil || cUser.ID == 0 {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	group.UserID = cUser.ID

	groupExists, err := group.NameExists()
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if groupExists {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Group already exists"})
	}

	lastInsertId, err := group.Create()
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Group created", Data: lastInsertId})
}

func (h *GroupHandler) GroupUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	var updateData map[string]any
	if err := config.ReadJSON(w, r, &updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	// get auth user in context
	cUser, ok := r.Context().Value(config.CKey("user")).(*models.User)

	if !ok || cUser == nil || cUser.ID == 0 {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	groups := &models.Group{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	groupExists, err := groups.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !groupExists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Group not found"})
	}

	queryParts := []string{"UPDATE groups SET"}
	params := []any{}
	paramCount := 1

	value, exists := updateData["name"]
	if exists {
		queryParts = append(queryParts, fmt.Sprintf("name=$%d,", paramCount))
		params = append(params, value)
		paramCount++
	}
	value, exists = updateData["uid"]
	if exists {
		queryParts = append(queryParts, fmt.Sprintf("uid=$%d,", paramCount))
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

	err = groups.Update(query, params)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (h *GroupHandler) GroupDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	// get auth user in context
	cUser, ok := r.Context().Value(config.CKey("user")).(*models.User)

	if !ok || cUser == nil || cUser.ID == 0 {
		return config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Credentials"})
	}

	groups := &models.Group{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	groupExists, err := groups.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !groupExists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Group not found"})
	}

	err = groups.Delete(id, cUser.ID)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Soft delte success"})
}
