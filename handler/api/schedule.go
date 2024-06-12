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
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
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
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Group not found"})
	}

	// request check
	request := &models.Request{}
	exists, err = request.IDExists(schedule.RequestID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
	}

	// notification check
	notification := &models.Notification{}
	exists, err = notification.IDExists(schedule.NotificationID, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Notification not found"})
	}

	// timinng check with request
	exists, err = schedule.TimingExists(cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusCreated, config.Response{Status: false, Message: "Timing already exists"})
	}

	err = schedule.Create(config.App().DB.DB)
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
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Schedule not found"})
	}

	queryParts := []string{"UPDATE schedules SET"}
	params := []any{}
	paramCount := 1

	if updateData.GroupID > 0 {
		// group check
		group := &models.Group{}
		exists, err := group.IDExists(schedule.GroupID, cUser.ID)
		if err != nil {
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
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
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
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
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
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
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
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
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}
	if !exists {
		return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Schedule not found"})
	}

	scheduleLogs, err := scheduleLog.Get(id, schedule_id, cUser.ID)
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusOK, config.Response{Status: true, Message: "Success", Data: scheduleLogs})
}

func (h *ScheduleHandler) ScheduleBulkHandler(w http.ResponseWriter, r *http.Request) error {
	bulk := &models.ScheduleBulk{}
	if err := config.ReadJSON(w, r, bulk); err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: err.Error()})
	}

	err := config.Validate(bulk)
	if err != nil {
		return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Content validation invalid", Data: err.Error()})
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	schedule := &models.Schedule{
		UserID:         cUser.ID,
		GroupID:        bulk.GroupID,
		RequestID:      bulk.RequestID,
		NotificationID: bulk.NotificationID,
		Timing:         bulk.Timing,
		Timeout:        bulk.Timeout,
		Retries:        bulk.Retries,
		Active:         bulk.Active,
	}

	tx, err := config.App().DB.Begin()
	if err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	// group check
	group := &models.Group{
		UserID: cUser.ID,
	}
	if schedule.GroupID > 0 {
		exists, err := group.IDExists(schedule.GroupID, cUser.ID)
		if err != nil {
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		if !exists {
			return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Group not found"})
		}
	} else {
		if bulk.Group == nil {
			return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Group or Group ID required"})
		}

		group.Name = bulk.Group.Name
		group.Active = bulk.Group.Active

		err = group.Create(tx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		schedule.GroupID = group.ID
	}

	// request check
	request := &models.Request{
		UserID: cUser.ID,
	}
	if schedule.RequestID > 0 {
		exists, err := request.IDExists(schedule.RequestID, cUser.ID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		if !exists {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Request not found"})
		}
	} else {
		if bulk.Request == nil {
			return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Request or Request ID required"})
		}

		request.Url = bulk.Request.Url
		request.Method = bulk.Request.Method
		request.Content = bulk.Request.Content
		request.Active = bulk.Request.Active

		exists, err := request.UrlExists()
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		if exists {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Url already exists"})
		}
		err = request.Create(tx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		for _, header := range bulk.Request.RequestHeaders {
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
		schedule.RequestID = request.ID
	}

	// notification check
	notification := &models.Notification{
		UserID: cUser.ID,
	}

	if schedule.NotificationID > 0 {
		exists, err := notification.IDExists(schedule.NotificationID, cUser.ID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		if !exists {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusNotFound, config.Response{Status: false, Message: "Notification not found"})
		}
	} else {
		if bulk.Notification == nil {
			return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Notification or Notification ID required"})
		}

		notification.Title = bulk.Notification.Title
		notification.Content = bulk.Notification.Content
		notification.IsMail = bulk.Notification.IsMail
		notification.IsSms = bulk.Notification.IsSms
		notification.Active = bulk.Notification.Active

		exists, err := notification.TitleExists()
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		if exists {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Title already exists"})
		}
		err = notification.Create(tx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
			}
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		for _, email := range bulk.Notification.NotifyEmails {
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

		for _, sms := range bulk.Notification.NotifySmses {
			notifySms := &models.NotifySms{
				NotificationID: notification.ID,
				Phone:          sms.Phone,
				Active:         sms.Active,
			}

			// check header key
			exists, err = notifySms.PhoneExists(tx, cUser.ID)
			if err != nil || exists {
				continue
			}

			err = notifySms.Create(tx)
			if err != nil {
				continue
			}
		}
		schedule.NotificationID = notification.ID
	}

	err = schedule.Create(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
		}
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return config.WriteJSON(w, http.StatusInternalServerError, config.Response{Status: false, Message: err.Error()})
	}

	return config.WriteJSON(w, http.StatusCreated, config.Response{Status: true, Message: "Schedule created", Data: schedule})
}
