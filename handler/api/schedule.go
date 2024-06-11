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

type ScheduleHandler struct{}

func (h *ScheduleHandler) ScheduleListHandler(w http.ResponseWriter, r *http.Request) error {
	schedule := &models.Schedule{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	group_id, _ := strconv.Atoi(r.URL.Query().Get("group_id"))
	request_id, _ := strconv.Atoi(r.URL.Query().Get("request_id"))
	notification_id, _ := strconv.Atoi(r.URL.Query().Get("notification_id"))
	timing := r.URL.Query().Get("timing")

	requests, err := schedule.Get(id, cUser.ID, group_id, request_id, notification_id, timing)
	if err != nil {
		return config.WriteJSON(w, http.StatusOK, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: requests})
}

func (h *ScheduleHandler) ScheduleCreateHandler(w http.ResponseWriter, r *http.Request) error {
	schedule := &models.Schedule{}
	if err := config.ReadJSON(w, r, schedule); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(schedule)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	schedule.UserID = cUser.ID

	// group check
	groups := &models.Group{}
	exists, err := groups.IDExists(schedule.GroupID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Group not found"})
	}

	// request check
	request := &models.Request{}
	exists, err = request.IDExists(schedule.RequestID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
	}

	// notification check
	notification := &models.Notification{}
	exists, err = notification.IDExists(schedule.NotificationID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Notification not found"})
	}

	err = schedule.Create()
	if err != nil {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Schedule created", Data: schedule})
}

func (h *ScheduleHandler) ScheduleUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	updateData := &models.ScheduleUpdate{}
	if err := config.ReadJSON(w, r, &updateData); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(updateData)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	schedule := &models.Schedule{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := schedule.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Schedule not found"})
	}

	queryParts := []string{"UPDATE schedules SET"}
	params := []any{}
	paramCount := 1

	if updateData.GroupID > 0 {
		// group check
		groups := &models.Group{}
		exists, err := groups.IDExists(schedule.GroupID, cUser.ID)
		if err != nil {
			return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
		}
		if !exists {
			return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Group not found"})
		}
		// add query
		queryParts = append(queryParts, fmt.Sprintf("group_id=$%d,", paramCount))
		params = append(params, updateData.GroupID)
		paramCount++
	}
	if updateData.RequestID > 0 {
		// request check
		request := &models.Request{}
		exists, err = request.IDExists(schedule.RequestID, cUser.ID)
		if err != nil {
			return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
		}
		if !exists {
			return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
		}
		// add query
		queryParts = append(queryParts, fmt.Sprintf("request_id=$%d,", paramCount))
		params = append(params, updateData.RequestID)
		paramCount++
	}
	if updateData.NotificationID > 0 {
		// request check
		notification := &models.Notification{}
		exists, err = notification.IDExists(schedule.NotificationID, cUser.ID)
		if err != nil {
			return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
		}
		if !exists {
			return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Notification not found"})
		}
		// add query
		queryParts = append(queryParts, fmt.Sprintf("notification_id=$%d,", paramCount))
		params = append(params, updateData.NotificationID)
		paramCount++
	}
	if updateData.Timing != "" {
		queryParts = append(queryParts, fmt.Sprintf("timing=$%d,", paramCount))
		params = append(params, updateData.Timing)
		paramCount++
	}
	if updateData.Timeout != nil {
		queryParts = append(queryParts, fmt.Sprintf("timeout=$%d,", paramCount))
		params = append(params, updateData.Timeout)
		paramCount++
	}
	if updateData.Retries != nil {
		queryParts = append(queryParts, fmt.Sprintf("retries=$%d,", paramCount))
		params = append(params, updateData.Retries)
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

	err = schedule.Update(query, params)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: updateData})
}

func (h *ScheduleHandler) ScheduleDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	schedule := &models.Schedule{}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	exists, err := schedule.IDExists(id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Schedule not found"})
	}

	err = schedule.Delete(id, cUser.ID)

	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Soft delte success"})
}

func (h *ScheduleHandler) ScheduleLogListHandler(w http.ResponseWriter, r *http.Request) error {
	scheduleLog := &models.ScheduleLog{}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	schedule_id, _ := strconv.Atoi(r.URL.Query().Get("schedule_id"))
	if schedule_id == 0 {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "schedule_id param required"})
	}

	schedule := &models.Schedule{}
	exists, err := schedule.IDExists(schedule_id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Schedule not found"})
	}

	scheduleLogs, err := scheduleLog.Get(id, schedule_id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusOK, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: scheduleLogs})
}
