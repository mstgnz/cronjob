package web

import (
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

type ScheduleHandler struct {
	schedule     *services.ScheduleService
	group        *services.GroupService
	request      *services.RequestService
	notification *services.NotificationService
}

func (h *ScheduleHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	_, group := h.group.ListService(w, r)
	_, requests := h.request.ListService(w, r)
	_, schedules := h.schedule.ListService(w, r)
	_, notifications := h.notification.ListService(w, r)
	return services.Render(w, r, "schedule", map[string]any{"lists": schedules.Data, "groups": group.Data, "requests": requests.Data, "notifications": notifications.Data}, "schedule/list", "schedule/log", "schedule/new")
}

func (h *ScheduleHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringIDsToInt(r, "group_id")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	jsonData, err = config.ConvertStringIDsToInt(r, "request_id")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	jsonData, err = config.ConvertStringIDsToInt(r, "notification_id")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	jsonData, err = config.ConvertStringIDsToInt(r, "timeout")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	jsonData, err = config.ConvertStringIDsToInt(r, "retries")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	jsonData, err = config.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.schedule.CreateService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/schedules")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *ScheduleHandler) EditHandler(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	query := r.URL.Query()
	query.Set("id", id)
	r.URL.RawQuery = query.Encode()

	_, response := h.schedule.ListService(w, r)

	data, _ := response.Data["schedules"].([]*models.Schedule)
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
		<tr hx-put="/schedules/%d" hx-trigger='cancel' hx-ext="json-enc" class='editing'>
			<th scope="row">%d</th>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td><input type="text" class="form-control" name="timing" value="%s"></td>
            <td><input type="number" class="form-control" name="timeout" value="%d"></td>
            <td><input type="number" class="form-control" name="retries" value="%d"></td>
            <td>%v</td>
            <td><select class="form-select" name="active">
				<option value="true" %s>Active</option>
				<option value="false" %s>Deactive</option>
			</select></td>
            <td>%s</td>
            <td>%s</td>
			<td>
				<div class="hstack gap-1">
				<button class="btn btn-warning" hx-get="/schedules">Cancel</button>
				<button class="btn btn-danger" hx-put="/schedules/%d" hx-include="closest tr">Save</button>
				</div>
			</td>
		</tr>
	`, data[0].ID, data[0].ID, data[0].Group.Name, data[0].Request.Url, data[0].Notification.Title, data[0].Timing, data[0].Timeout, data[0].Retries, data[0].Running, activeSelected, deactiveSelected, data[0].CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, data[0].ID)

	_, _ = w.Write([]byte(form))
	return nil
}

func (h *ScheduleHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringIDsToInt(r, "timeout")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	jsonData, err = config.ConvertStringIDsToInt(r, "retries")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	jsonData, err = config.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.schedule.UpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/schedules")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *ScheduleHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.schedule.DeleteService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/schedules")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *ScheduleHandler) PaginationHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	schedule := &models.Schedule{}

	search := ""
	total := schedule.Count(cUser.ID)
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/schedules-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Previous</button>
        </li>`, previous)

		for i := 1; i <= size; i++ {
			pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/schedules-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">%d</button>
        </li>`, i, i)
		}

		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/schedules-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Next</button>
        </li>`, next)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	schedules := schedule.WithQuery(cUser.ID, (current-1)*row, row, search)

	tr := ""
	for _, v := range schedules {
		if v.DeletedAt != nil {
			continue
		}
		var updatedAt = ""
		if v.UpdatedAt != nil {
			updatedAt = v.UpdatedAt.Format("2006-01-02 15:04:05")
		}

		tr += fmt.Sprintf(`<tr>
            <th scope="row">%d</th>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%d</td>
            <td>%d</td>
            <td>%v</td>
            <td>%v</td>
            <td>%s</td>
            <td>%s</td>
            <td>
				<div class="hstack gap-1">
					<button class="btn btn-info" hx-get="/schedules/%d"
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
					<button class="btn btn-danger" hx-delete="/schedules/%d"  hx-trigger='confirmed' onClick="Swal.fire({
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
        </tr>`, v.ID, v.Group.Name, v.Request.Url, v.Notification.Title, v.Timing, v.Timeout, v.Retries, v.Running, v.Active, v.CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, v.ID, v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}

func (h *ScheduleHandler) LogPaginationHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	scheduleLog := &models.ScheduleLog{}

	search := ""
	total := scheduleLog.Count(cUser.ID)
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/schedules/logs?page=%d" hx-target="#log" hx-swap="innerHTML" hx-trigger="click">Previous</button>
        </li>`, previous)

		for i := 1; i <= size; i++ {
			pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/schedules/logs?page=%d" hx-target="#log" hx-swap="innerHTML" hx-trigger="click">%d</button>
        </li>`, i, i)
		}

		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/schedules/logs?page=%d" hx-target="#log" hx-swap="innerHTML" hx-trigger="click">Next</button>
        </li>`, next)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	schedules := scheduleLog.Paginate(cUser.ID, (current-1)*row, row, search)

	tr := ""
	for _, v := range schedules {
		tr += fmt.Sprintf(`<tr>
            <th scope="row">%d</th>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%f</td>
            <td>%v</td>
            <td>%s</td>
        </tr>`, v.ID, v.Schedule.Timing, v.StartedAt.Format("2006-01-02 15:04:05"), v.FinishedAt.Format("2006-01-02 15:04:05"), v.Took, v.Result, v.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	_, _ = w.Write([]byte(tr))
	return nil
}
