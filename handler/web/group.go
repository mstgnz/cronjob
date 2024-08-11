package web

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/pkg/config"
	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/services"
)

type GroupHandler struct {
	*services.GroupService
}

func (h *GroupHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	_, data := h.ListService(w, r)
	return services.Render(w, r, "group", map[string]any{"lists": data.Data}, "group/list", "group/new")
}

func (h *GroupHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {

	jsonData, err := response.ConvertStringIDsToInt(r, "uid")
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

	code, response := h.CreateService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/groups")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *GroupHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := response.ConvertStringIDsToInt(r, "uid")
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

	code, response := h.UpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/groups")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *GroupHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.DeleteService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/groups")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *GroupHandler) PaginationHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	group := &models.Group{}

	search := ""
	total := group.Count(cUser.ID)
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/groups-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Previous</button>
        </li>`, previous)

		for i := 1; i <= size; i++ {
			pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/groups-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">%d</button>
        </li>`, i, i)
		}

		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/groups-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Next</button>
        </li>`, next)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	groups := group.Paginate(cUser.ID, (current-1)*row, row, search)

	tr := ""
	for _, v := range groups {
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
            <td>%v</td>
            <td>%s</td>
            <td>%s</td>
            <td>
				<div class="hstack gap-1">
					<button class="btn btn-info" hx-get="/groups/%d"
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
					<button class="btn btn-danger" hx-delete="/groups/%d"  hx-trigger='confirmed' onClick="Swal.fire({
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
        </tr>`, v.ID, v.Name, v.Parent.Name, v.Active, v.CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, v.ID, v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}

func (h *GroupHandler) EditHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	id := chi.URLParam(r, "id")
	query := r.URL.Query()
	query.Set("id", id)
	r.URL.RawQuery = query.Encode()

	_, response := h.ListService(w, r)

	data, _ := response.Data["groups"].([]*models.Group)
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

	group := &models.Group{}
	groups, _ := group.Get(cUser.ID, 0, 0)

	groupSelect := `<select name="uid" class="form-control">`
	for _, v := range groups {
		groupSelect += fmt.Sprintf(`<option value="%d">%s</option>`, v.ID, v.Name)
	}
	groupSelect += `</select>`

	form := fmt.Sprintf(`
		<tr hx-put="/groups/%d" hx-trigger='cancel' hx-target="#toast" hx-swap="innerHTML" hx-ext="json-enc" class='editing'>
			<th scope="row">%d</th>
            <td><input name="name" class="form-control" value="%s"></td>
            <td>%v</td>
            <td><select class="form-select" name="active">
                    <option value="true" %s>Active</option>
                    <option value="false" %s>Deactive</option>
                </select></td>
            <td>%s</td>
            <td>%s</td>
			<td>
				<div class="hstack gap-1">
				<button class="btn btn-warning" hx-get="/groups-pagination" hx-target="#tbody" hx-swap="innerHTML">Cancel</button>
				<button class="btn btn-danger" hx-put="/groups/%d" hx-include="closest tr">Save</button>
				</div>
			</td>
		</tr>
	`, data[0].ID, data[0].ID, data[0].Name, groupSelect, activeSelected, deactiveSelected, data[0].CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, data[0].ID)

	_, _ = w.Write([]byte(form))
	return nil
}
