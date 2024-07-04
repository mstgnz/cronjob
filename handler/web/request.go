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

type RequestHandler struct {
	request *services.RequestService
	header  *services.RequestHeaderService
}

func (h *RequestHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	_, requests := h.request.ListService(w, r)
	return services.Render(w, r, "request", map[string]any{"lists": requests.Data}, "request/list", "request/header", "request/new")
}

func (h *RequestHandler) CreateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.request.CreateService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/requests")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *RequestHandler) EditHandler(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	query := r.URL.Query()
	query.Set("id", id)
	r.URL.RawQuery = query.Encode()

	_, response := h.request.ListService(w, r)

	data, _ := response.Data["requests"].([]models.Request)
	var updatedAt = ""
	if data[0].UpdatedAt != nil {
		updatedAt = data[0].UpdatedAt.Format("2006-01-02 15:04:05")
	}

	methodGetSelected := ""
	methodPostSelected := ""
	methodPutSelected := ""
	methodPatchSelected := ""
	methodDeleteSelected := ""
	activeSelected := ""
	deactiveSelected := ""

	switch data[0].Method {
	case "GET":
		methodGetSelected = "selected"
	case "POST":
		methodPostSelected = "selected"
	case "PUT":
		methodPutSelected = "selected"
	case "PATCH":
		methodPatchSelected = "selected"
	case "DELETE":
		methodDeleteSelected = "selected"
	}

	if data[0].Active {
		activeSelected = "selected"
	} else {
		deactiveSelected = "selected"
	}

	form := fmt.Sprintf(`
		<tr hx-put="/requests/%d" hx-trigger='cancel'  hx-ext="json-enc" class='editing'>
			<th scope="row">%d</th>
            <td><input name="url" class="form-control" value="%s"></td>
            <td><select class="form-select" name="method">
                    <option value="GET" %s>GET</option>
                    <option value="POST" %s>POST</option>
                    <option value="PUT" %s>PUT</option>
                    <option value="PATCH" %s>PATCH</option>
                    <option value="DELETE" %s>DELETE</option>
                </select></td>
            <td><textarea class="form-control" name="content" placeholder="Content">%s</textarea></td>
            <td><select class="form-select" name="active">
                    <option value="true" %s>Active</option>
                    <option value="false" %s>Deactive</option>
                </select></td>
            <td>%s</td>
            <td>%s</td>
			<td>
				<div class="hstack gap-1">
				<button class="btn btn-warning" hx-get="/requests">Cancel</button>
				<button class="btn btn-danger" hx-put="/requests/%d" hx-include="closest tr">Save</button>
				</div>
			</td>
		</tr>
	`, data[0].ID, data[0].ID, data[0].Url, methodGetSelected, methodPostSelected, methodPutSelected, methodPatchSelected, methodDeleteSelected, data[0].Content, activeSelected, deactiveSelected, data[0].CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, data[0].ID)

	_, _ = w.Write([]byte(form))
	return nil
}

func (h *RequestHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.request.UpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/requests")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *RequestHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.request.DeleteService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/requests")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *RequestHandler) PaginationHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	request := &models.Request{}

	search := ""
	total := request.Count(cUser.ID)
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/requests-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Previous</button>
        </li>`, previous)

		for i := 1; i <= size; i++ {
			pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/requests-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">%d</button>
        </li>`, i, i)
		}

		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/requests-pagination?page=%d" hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click">Next</button>
        </li>`, next)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	requests := request.Paginate(cUser.ID, (current-1)*row, row, search)

	tr := ""
	for _, v := range requests {
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
            <td>%s</td>
            <td>%v</td>
            <td>%s</td>
            <td>%s</td>
            <td>
				<div class="hstack gap-1">
					<button class="btn btn-info" data-request='%s' hx-get="/requests/%d/edit"
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
					<button class="btn btn-danger" hx-delete="/requests/%d"  hx-trigger='confirmed' onClick="Swal.fire({
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
        </tr>`, v.ID, v.Url, v.Method, v.Content, v.Active, v.CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, dataRequest, v.ID, v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}

func (h *RequestHandler) HeaderCreateHandler(w http.ResponseWriter, r *http.Request) error {

	jsonData, err := config.ConvertStringIDsToInt(r, "request_id")
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

	code, response := h.header.CreateService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/requests")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *RequestHandler) HeaderUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringBoolsToBool(r, "active")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}
	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.header.UpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/requests")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *RequestHandler) HeaderDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.header.DeleteService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/requests")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *RequestHandler) HeaderPaginationHandler(w http.ResponseWriter, r *http.Request) error {
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	requestHeader := &models.RequestHeader{}

	search := ""
	total := requestHeader.Count(cUser.ID)
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/requests-pagination?page=%d" hx-target="#header" hx-swap="innerHTML" hx-trigger="click">Previous</button>
        </li>`, previous)

		for i := 1; i <= size; i++ {
			pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/requests-pagination?page=%d" hx-target="#header" hx-swap="innerHTML" hx-trigger="click">%d</button>
        </li>`, i, i)
		}

		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/requests-pagination?page=%d" hx-target="#header" hx-swap="innerHTML" hx-trigger="click">Next</button>
        </li>`, next)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	requestHeaders := requestHeader.Paginate(cUser.ID, (current-1)*row, row, search)

	tr := ""
	for _, v := range requestHeaders {
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
					<button class="btn btn-info" data-request='%s' hx-get="/requests/headers/%d/edit"
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
					<button class="btn btn-danger" hx-delete="/requests/headers/%d" hx-trigger='confirmed' onClick="Swal.fire({
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
        </tr>`, v.ID, v.Key, v.Value, v.Active, v.CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, dataHeader, v.ID, v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}

func (h *RequestHandler) HeaderEditHandler(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	query := r.URL.Query()
	query.Set("id", id)
	r.URL.RawQuery = query.Encode()

	_, response := h.header.ListService(w, r)

	data, _ := response.Data["request_headers"].([]models.RequestHeader)
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
		<tr hx-put="/requests/headers/%d" hx-trigger='cancel'  hx-ext="json-enc" class='editing'>
			<th scope="row">%d</th>
            <td><input name="key" class="form-control" value="%s" /></td>
            <td><input name="value" class="form-control" value="%s" /></td>
            <td><select class="form-select" name="active">
                    <option value="true" %s>Active</option>
                    <option value="false" %s>Deactive</option>
                </select></td>
            <td>%s</td>
            <td>%s</td>
			<td>
				<div class="hstack gap-1">
				<button class="btn btn-warning" hx-get="/requests">Cancel</button>
				<button class="btn btn-danger" hx-put="/requests/headers/%d" hx-include="closest tr">Save</button>
				</div>
			</td>
		</tr>
	`, data[0].ID, data[0].ID, data[0].Key, data[0].Value, activeSelected, deactiveSelected, data[0].CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, data[0].ID)

	_, _ = w.Write([]byte(form))
	return nil
}
