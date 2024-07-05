package web

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/services"
)

type WebhookHandler struct {
	webhook  *services.WebhookService
	request  *services.RequestService
	schedule *services.ScheduleService
}

func (h *WebhookHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return services.Render(w, r, "webhook", map[string]any{}, "webhook/list", "webhook/new")
}

func (h *WebhookHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.webhook.CreateService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/webhooks")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *WebhookHandler) EditHandler(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	query := r.URL.Query()
	query.Set("id", id)
	r.URL.RawQuery = query.Encode()

	_, response := h.webhook.ListService(w, r)

	data, _ := response.Data["webhooks"].([]*models.Webhook)
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

	// requests
	_, requests := h.request.ListService(w, r)
	requestData, _ := requests.Data["requests"].([]models.Request)
	requestList := `<select class="form-select" name="request_id">`
	for _, request := range requestData {
		requestList += fmt.Sprintf(`<option value="%d" %s>GET</option>`, request.ID, request.Url)
	}
	requestList += "</select>"

	// schedules
	_, schedules := h.schedule.ListService(w, r)
	schedulesData, _ := schedules.Data["schedules"].([]models.Schedule)
	scheduleList := `<select class="form-select" name="schedule_id">`
	for _, schedule := range schedulesData {
		scheduleList += fmt.Sprintf(`<option value="%d" %s>GET</option>`, schedule.ID, schedule.Timing)
	}
	scheduleList += "</select>"

	form := fmt.Sprintf(`
		<tr hx-put="/webhooks/%d" hx-trigger='cancel'  hx-ext="json-enc" class='editing'>
			<th scope="row">%d</th>
            <td>%s</td>
            <td>%s</td>
            <td><select class="form-select" name="active">
                    <option value="true" %s>Active</option>
                    <option value="false" %s>Deactive</option>
                </select></td>
			<td>%s</td>
			<td>%s</td>
			<td>
				<div class="hstack gap-1">
				<button class="btn btn-warning" hx-get="/webhooks">Cancel</button>
				<button class="btn btn-danger" hx-put="/webhooks/%d" hx-include="closest tr">Save</button>
				</div>
			</td>
		</tr>
	`, data[0].ID, data[0].ID, scheduleList, requestList, activeSelected, deactiveSelected, data[0].CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, data[0].ID)

	_, _ = w.Write([]byte(form))
	return nil
}

func (h *WebhookHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.webhook.UpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/webhooks")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *WebhookHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.webhook.DeleteService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/webhooks")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *WebhookHandler) PaginationHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	webhook := &models.Webhook{}

	search := ""
	total := webhook.Count(cUser.ID)
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/webhooks-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Previous</button>
        </li>`, previous)

		for i := 1; i <= size; i++ {
			pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/webhooks-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">%d</button>
        </li>`, i, i)
		}

		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/webhooks-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Next</button>
        </li>`, next)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	webhooks := webhook.Paginate(cUser.ID, (current-1)*row, row, search)

	tr := ""
	for _, v := range webhooks {
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
            <td>%s</td>
            <td>%s</td>
            <td>
				<div class="hstack gap-1">
					<button class="btn btn-info" data-request='%s' hx-get="/webhooks/%d"
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
					<button class="btn btn-danger" hx-delete="/webhooks/%d"  hx-trigger='confirmed' onClick="Swal.fire({
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
        </tr>`, v.ID, v.Schedule.Timing, v.Request.Url, v.Active, v.CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, dataRequest, v.ID, v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}
