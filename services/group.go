package services

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

type GroupService struct{}

func (s *GroupService) ListSerice(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	group := &models.Group{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))

	groups, err := group.Get(cUser.ID, id, uid)
	if err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, config.Response{Status: true, Message: "Success", Data: groups}
}

func (s *GroupService) CreateSerice(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	group := &models.Group{}
	if err := config.ReadJSON(w, r, group); err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: err.Error()}
	}

	err := config.Validate(group)
	if err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	group.UserID = cUser.ID

	exists, err := group.NameExists()
	if err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}
	if exists {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Group already exists"}
	}

	err = group.Create(config.App().DB.DB)
	if err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}

	return http.StatusCreated, config.Response{Status: true, Message: "Group created", Data: group}
}

func (s *GroupService) UpdateSerice(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	updateData := &models.GroupUpdate{}
	if err := config.ReadJSON(w, r, &updateData); err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: err.Error()}
	}

	err := config.Validate(updateData)
	if err != nil {
		return http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	groups := &models.Group{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := groups.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, config.Response{Status: false, Message: "Group not found"}
	}

	queryParts := []string{"UPDATE groups SET"}
	params := []any{}
	paramCount := 1

	if updateData.Name != "" {
		queryParts = append(queryParts, fmt.Sprintf("name=$%d,", paramCount))
		params = append(params, updateData.Name)
		paramCount++
	}
	if updateData.UID > 0 {
		queryParts = append(queryParts, fmt.Sprintf("uid=$%d,", paramCount))
		params = append(params, updateData.UID)
		paramCount++
	}
	if updateData.Active != nil {
		queryParts = append(queryParts, fmt.Sprintf("active=$%d,", paramCount))
		params = append(params, updateData.Active)
		paramCount++
	}

	if len(params) == 0 {
		return http.StatusBadRequest, config.Response{Status: false, Message: "No fields to update"}
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
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData}
}

func (s *GroupService) DeleteSerice(w http.ResponseWriter, r *http.Request) (int, config.Response) {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	groups := &models.Group{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := groups.IDExists(id, cUser.ID)
	if err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}
	if !exists {
		return http.StatusNotFound, config.Response{Status: false, Message: "Group not found"}
	}

	err = groups.Delete(id, cUser.ID)

	if err != nil {
		return http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()}
	}

	return http.StatusOK, config.Response{Status: true, Message: "Soft delte success"}
}
