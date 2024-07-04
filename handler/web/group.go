package web

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
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

	jsonData, err := config.ConvertStringIDsToInt(r, "uid")
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

	code, response := h.CreateService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/groups")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *GroupHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringIDsToInt(r, "uid")
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
		dataGroup, _ := json.Marshal(v)
		tr += fmt.Sprintf(`<tr>
            <th scope="row">%d</th>
            <td>%s</td>
            <td>%s</td>
            <td>%v</td>
            <td>%s</td>
            <td>%s</td>
            <td>
				<div class="hstack gap-1">
					<button class="btn btn-info" data-group='%s'>
						<i class="bi bi-pencil"></i>
					</button>
					<button class="btn btn-danger" hx-delete="/groups/%d" hx-confirm="Are you sure?">
						<i class="bi bi-trash-fill"></i>
					</button>
				</div>
			</td>
        </tr>`, v.ID, v.Name, v.Parent.Name, v.Active, v.CreatedAt.Format("2006-01-02 15:04:05"), updatedAt, dataGroup, v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}
