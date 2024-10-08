package web

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/pkg/config"
	"github.com/mstgnz/cronjob/pkg/load"
	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/services"
)

type NotificationHandler struct {
	notify  *services.NotificationService
	email   *services.NotifyEmailService
	message *services.NotifyMessageService
}

func (h *NotificationHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	_, notification := h.notify.ListService(w, r)
	return load.Render(w, r, "notification", map[string]any{"lists": notification.Data}, "notification/list", "notification/email", "notification/message", "notification/new")
}

func (h *NotificationHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := response.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))
	jsonData, err = response.ConvertStringBoolsToBool(r, "is_mail")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))
	jsonData, err = response.ConvertStringBoolsToBool(r, "is_message")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.notify.CreateService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/notifications")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *NotificationHandler) EditHandler(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	query := r.URL.Query()
	query.Set("id", id)
	r.URL.RawQuery = query.Encode()

	_, response := h.notify.ListService(w, r)

	data, _ := response.Data["notifications"].([]*models.Notification)
	var updatedAt = ""
	if data[0].UpdatedAt != nil {
		updatedAt = data[0].UpdatedAt.Format("2006-01-02 15:04:05")
	}

	activeSelected := ""
	deactiveSelected := ""
	activeEmailSelected := ""
	deactiveEmailSelected := ""
	activeMessageSelected := ""
	deactiveMessageSelected := ""

	if data[0].Active {
		activeSelected = "selected"
	} else {
		deactiveSelected = "selected"
	}
	if data[0].IsMail {
		activeEmailSelected = "selected"
	} else {
		deactiveEmailSelected = "selected"
	}
	if data[0].IsMessage {
		activeMessageSelected = "selected"
	} else {
		deactiveMessageSelected = "selected"
	}

	form := fmt.Sprintf(`
		<tr hx-put="/notifications/%d" hx-trigger='cancel' hx-target="#toast" hx-swap="innerHTML" hx-ext="json-enc" class='editing'>
			<th scope="row">%d</th>
            <td><input name="title" class="form-control" value="%s"></td>
            <td><input name="content" class="form-control" value="%s"></td>
            <td><select class="form-select" name="is_mail">
                    <option value="true" %s>Active</option>
                    <option value="false" %s>Deactive</option>
                </select></td>
			<td><select class="form-select" name="is_message">
                    <option value="true" %s>Active</option>
                    <option value="false" %s>Deactive</option>
                </select></td>
			<td><select class="form-select" name="active">
                    <option value="true" %s>Active</option>
                    <option value="false" %s>Deactive</option>
                </select></td>
            <td>%s</td>
            <td>%s</td>
			<td>
				<div class="hstack gap-1">
				<button class="btn btn-warning" hx-get="/notifications-pagination" hx-target="#tbody" hx-swap="innerHTML">Cancel</button>
				<button class="btn btn-danger" hx-put="/notifications/%d" hx-include="closest tr">Save</button>
				</div>
			</td>
		</tr>
	`, data[0].ID, data[0].ID, data[0].Title, data[0].Content, activeEmailSelected, deactiveEmailSelected, activeMessageSelected, deactiveMessageSelected, activeSelected, deactiveSelected, data[0].CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, data[0].ID)

	_, _ = w.Write([]byte(form))
	return nil
}

func (h *NotificationHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := response.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))
	jsonData, err = response.ConvertStringBoolsToBool(r, "is_mail")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))
	jsonData, err = response.ConvertStringBoolsToBool(r, "is_message")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.notify.UpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/notifications")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *NotificationHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.notify.DeleteService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/notifications")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *NotificationHandler) PaginationHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notification := &models.Notification{}

	search := ""
	total := notification.Count(cUser.ID)
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/notifications-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Previous</button>
        </li>`, previous)

		for i := 1; i <= size; i++ {
			pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/notifications-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">%d</button>
        </li>`, i, i)
		}

		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/notifications-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Next</button>
        </li>`, next)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	notifications := notification.Paginate(cUser.ID, (current-1)*row, row, search)

	tr := ""
	for _, v := range notifications {
		if v.DeletedAt != nil {
			continue
		}
		var updatedAt = ""
		if v.UpdatedAt != nil {
			updatedAt = v.UpdatedAt.Format("2006-01-02 15:04:05")
		}
		dataRequest, _ := json.Marshal(v)
		tr += fmt.Sprintf(`<tr>
            <th scope="row">%d</th>
            <td>%s</td>
            <td>%s</td>
            <td>%v</td>
            <td>%v</td>
            <td>%v</td>
            <td>%s</td>
            <td>%s</td>
            <td>
				<div class="hstack gap-1">
					<button class="btn btn-info" data-notifications='%s' hx-get="/notifications/%d"
						hx-trigger="edit"
						onClick="let editing = document.querySelector('.editing')
						if(editing) {
						Swal.fire({title: 'Already Editing',
									showCancelButton: true,
									confirmButtonText: 'Yep, Edit This Row!',
									text:'Hey!  You are already editing a row!  Do you want to cancel that edit and continue?'})
						.then((result) => {
							if(result.isConfirmed) {
								htmx.trigger(editing, 'cancel')
								htmx.trigger(this, 'edit')
							}
						})
						} else {
						htmx.trigger(this, 'edit')
						}">
						<i class="bi bi-pencil"></i>
					</button>
					<button class="btn btn-danger" hx-delete="/notifications/%d"  hx-trigger='confirmed' onClick="Swal.fire({
							title: 'Do you approve the deletion?',
							icon: 'warning',
							showCancelButton: true,
							cancelButtonColor: '#d33',
							cancelButtonText: 'Close',
							confirmButtonColor: '#3085d6',
							confirmButtonText: 'Yes Delete'
						}).then((result) => {if (result.isConfirmed) {htmx.trigger(this, 'confirmed')}})">
						<i class="bi bi-trash-fill"></i>	
					</button>
				</div>
			</td>
        </tr>`, v.ID, v.Title, v.Content, v.IsMail, v.IsMessage, v.Active, v.CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, dataRequest, v.ID, v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}

func (h *NotificationHandler) EmailCreateHandler(w http.ResponseWriter, r *http.Request) error {

	jsonData, err := response.ConvertStringIDsToInt(r, "notification_id")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	jsonData, err = response.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.email.CreateService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/notifications")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *NotificationHandler) EmailUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := response.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.email.UpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/notifications")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *NotificationHandler) EmailDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.email.DeleteService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/notifications")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *NotificationHandler) EmailPaginationHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notifyEmail := &models.NotifyEmail{}

	search := ""
	total := notifyEmail.Count(cUser.ID)
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/notifications/email-pagination?page=%d" hx-target="#email" hx-swap="innerHTML" hx-trigger="click">Previous</button>
        </li>`, previous)

		for i := 1; i <= size; i++ {
			pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/notifications/email-pagination?page=%d" hx-target="#email" hx-swap="innerHTML" hx-trigger="click">%d</button>
        </li>`, i, i)
		}

		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/notifications/email-pagination?page=%d" hx-target="#email" hx-swap="innerHTML" hx-trigger="click">Next</button>
        </li>`, next)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	notifyEmails := notifyEmail.Paginate(cUser.ID, (current-1)*row, row, search)

	tr := ""
	for _, v := range notifyEmails {
		if v.DeletedAt != nil {
			continue
		}
		var updatedAt = ""
		if v.UpdatedAt != nil {
			updatedAt = v.UpdatedAt.Format("2006-01-02 15:04:05")
		}
		dataHeader, _ := json.Marshal(v)
		tr += fmt.Sprintf(`<tr>
            <th scope="row">%d</th>
            <td>%s</td>
            <td>%s</td>
            <td>%v</td>
            <td>%s</td>
            <td>%s</td>
            <td>
				<div class="hstack gap-1">
					<button class="btn btn-info" data-email='%s' hx-get="/notifications/email/%d"
						hx-trigger="edit"
						onClick="let editing = document.querySelector('.editing')
						if(editing) {
						Swal.fire({title: 'Already Editing',
									showCancelButton: true,
									confirmButtonText: 'Yep, Edit This Row!',
									text:'Hey!  You are already editing a row!  Do you want to cancel that edit and continue?'})
						.then((result) => {
							if(result.isConfirmed) {
								htmx.trigger(editing, 'cancel')
								htmx.trigger(this, 'edit')
							}
						})
						} else {
						htmx.trigger(this, 'edit')
						}">
						<i class="bi bi-pencil"></i>
					</button>
					<button class="btn btn-danger" hx-delete="/notifications/email/%d" hx-trigger='confirmed' onClick="Swal.fire({
							title: 'Do you approve the deletion?',
							icon: 'warning',
							showCancelButton: true,
							cancelButtonColor: '#d33',
							cancelButtonText: 'Close',
							confirmButtonColor: '#3085d6',
							confirmButtonText: 'Yes Delete'
						}).then((result) => {if (result.isConfirmed) {htmx.trigger(this, 'confirmed')}})">
						<i class="bi bi-trash-fill"></i>
					</button>
				</div>
			</td>
        </tr>`, v.ID, v.Notification.Title, v.Email, v.Active, v.CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, dataHeader, v.ID, v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}

func (h *NotificationHandler) EmailEditHandler(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	query := r.URL.Query()
	query.Set("id", id)
	r.URL.RawQuery = query.Encode()

	_, response := h.email.ListService(w, r)

	data, _ := response.Data["notify_emails"].([]*models.NotifyEmail)
	var updatedAt = ""
	if data[0].UpdatedAt != nil {
		updatedAt = data[0].UpdatedAt.Format("2006-01-02 15:04:05")
	}

	activeSelected := ""
	deactiveSelected := ""
	if data[0].Active {
		activeSelected = "selected"
	} else {
		deactiveSelected = "selected"
	}

	form := fmt.Sprintf(`
		<tr hx-put="/notifications/email/%d" hx-trigger='cancel' hx-target="#toast" hx-swap="innerHTML" hx-ext="json-enc" class='editing'>
			<th scope="row">%d</th>
            <td>%s</td>
            <td><input name="email" class="form-control" value="%s" /></td>
            <td><select class="form-select" name="active">
                    <option value="true" %s>Active</option>
                    <option value="false" %s>Deactive</option>
                </select></td>
            <td>%s</td>
            <td>%s</td>
			<td>
				<div class="hstack gap-1">
				<button class="btn btn-warning" hx-get="/notifications/email-pagination" hx-target="#email" hx-swap="innerHTML">Cancel</button>
				<button class="btn btn-danger" hx-put="/notifications/email/%d" hx-include="closest tr">Save</button>
				</div>
			</td>
		</tr>
	`, data[0].ID, data[0].ID, data[0].Notification.Title, data[0].Email, activeSelected, deactiveSelected, data[0].CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, data[0].ID)

	_, _ = w.Write([]byte(form))
	return nil
}

func (h *NotificationHandler) MessageCreateHandler(w http.ResponseWriter, r *http.Request) error {

	jsonData, err := response.ConvertStringIDsToInt(r, "notification_id")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	jsonData, err = response.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.message.CreateService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/notifications")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *NotificationHandler) MessageUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := response.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.message.UpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/notifications")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *NotificationHandler) MessageDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.message.DeleteService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/notifications")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *NotificationHandler) MessagePaginationHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	notifyMessage := &models.NotifyMessage{}

	search := ""
	total := notifyMessage.Count(cUser.ID)
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/notifications/message-pagination?page=%d" hx-target="#message" hx-swap="innerHTML" hx-trigger="click">Previous</button>
        </li>`, previous)

		for i := 1; i <= size; i++ {
			pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/notifications/message-pagination?page=%d" hx-target="#message" hx-swap="innerHTML" hx-trigger="click">%d</button>
        </li>`, i, i)
		}

		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/notifications/message-pagination?page=%d" hx-target="#message" hx-swap="innerHTML" hx-trigger="click">Next</button>
        </li>`, next)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	notifyMessages := notifyMessage.Paginate(cUser.ID, (current-1)*row, row, search)

	tr := ""
	for _, v := range notifyMessages {
		if v.DeletedAt != nil {
			continue
		}
		var updatedAt = ""
		if v.UpdatedAt != nil {
			updatedAt = v.UpdatedAt.Format("2006-01-02 15:04:05")
		}
		dataHeader, _ := json.Marshal(v)
		tr += fmt.Sprintf(`<tr>
            <th scope="row">%d</th>
            <td>%s</td>
            <td>%s</td>
            <td>%v</td>
            <td>%s</td>
            <td>%s</td>
            <td>
				<div class="hstack gap-1">
					<button class="btn btn-info" data-message='%s' hx-get="/notifications/message/%d"
						hx-trigger="edit"
						onClick="let editing = document.querySelector('.editing')
						if(editing) {
						Swal.fire({title: 'Already Editing',
									showCancelButton: true,
									confirmButtonText: 'Yep, Edit This Row!',
									text:'Hey!  You are already editing a row!  Do you want to cancel that edit and continue?'})
						.then((result) => {
							if(result.isConfirmed) {
								htmx.trigger(editing, 'cancel')
								htmx.trigger(this, 'edit')
							}
						})
						} else {
						htmx.trigger(this, 'edit')
						}">
						<i class="bi bi-pencil"></i>
					</button>
					<button class="btn btn-danger" hx-delete="/notifications/message/%d" hx-trigger='confirmed' onClick="Swal.fire({
							title: 'Do you approve the deletion?',
							icon: 'warning',
							showCancelButton: true,
							cancelButtonColor: '#d33',
							cancelButtonText: 'Close',
							confirmButtonColor: '#3085d6',
							confirmButtonText: 'Yes Delete'
						}).then((result) => {if (result.isConfirmed) {htmx.trigger(this, 'confirmed')}})">
						<i class="bi bi-trash-fill"></i>
					</button>
				</div>
			</td>
        </tr>`, v.ID, v.Notification.Title, v.Phone, v.Active, v.CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, dataHeader, v.ID, v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}

func (h *NotificationHandler) MessageEditHandler(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	query := r.URL.Query()
	query.Set("id", id)
	r.URL.RawQuery = query.Encode()

	_, response := h.message.ListService(w, r)

	data, _ := response.Data["notify_messages"].([]*models.NotifyMessage)
	var updatedAt = ""
	if data[0].UpdatedAt != nil {
		updatedAt = data[0].UpdatedAt.Format("2006-01-02 15:04:05")
	}

	activeSelected := ""
	deactiveSelected := ""
	if data[0].Active {
		activeSelected = "selected"
	} else {
		deactiveSelected = "selected"
	}

	form := fmt.Sprintf(`
		<tr hx-put="/notifications/message/%d" hx-trigger='cancel' hx-target="#toast" hx-swap="innerHTML" hx-ext="json-enc" class='editing'>
			<th scope="row">%d</th>
            <td>%s</td>
            <td><input name="text" class="form-control" value="%s" /></td>
            <td><select class="form-select" name="active">
                    <option value="true" %s>Active</option>
                    <option value="false" %s>Deactive</option>
                </select></td>
            <td>%s</td>
            <td>%s</td>
			<td>
				<div class="hstack gap-1">
				<button class="btn btn-warning" hx-get="/notifications/message-pagination" hx-target="#message" hx-swap="innerHTML">Cancel</button>
				<button class="btn btn-danger" hx-put="/notifications/message/%d" hx-include="closest tr">Save</button>
				</div>
			</td>
		</tr>
	`, data[0].ID, data[0].ID, data[0].Notification.Title, data[0].Phone, activeSelected, deactiveSelected, data[0].CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, data[0].ID)

	_, _ = w.Write([]byte(form))
	return nil
}
