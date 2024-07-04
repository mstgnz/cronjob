package web

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/services"
)

type SettingHandler struct {
	user *services.UserService
}

func (h *SettingHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return services.Render(w, r, "setting", map[string]any{}, "setting/user-list", "setting/user-new", "setting/app-log")
}

func (h *SettingHandler) UserSignUpHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.user.RegisterService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/setting")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *SettingHandler) UsersHandler(w http.ResponseWriter, r *http.Request) error {
	user := &models.User{}

	search := ""
	total := user.Count()
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	users := user.Get((current-1)*row, row, search)

	data := map[string]any{}
	data["previous"] = previous
	data["next"] = next
	data["current"] = current
	data["size"] = size
	data["users"] = users

	w.Header().Set("Content-Type", "application/json")
	result := config.MaptoJSON(data)
	_, _ = w.Write(result)
	return nil
}

func (h *SettingHandler) UserChangeProfileHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringIDsToInt(r, "id")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}

	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.user.UpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/setting")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *SettingHandler) UserChangePasswordHandler(w http.ResponseWriter, r *http.Request) error {
	jsonData, err := config.ConvertStringIDsToInt(r, "id")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return nil
	}

	r.Body = io.NopCloser(strings.NewReader(string(jsonData)))

	code, response := h.user.PassUpdateService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/setting")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *SettingHandler) UserDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.user.DeleteService(w, r)
	if response.Status && code == http.StatusOK {
		w.Header().Set("HX-Redirect", "/setting")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *SettingHandler) AppLogHandler(w http.ResponseWriter, r *http.Request) error {
	appLog := &models.AppLog{}

	search := ""
	total := appLog.Count()
	row := 20

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	activeClass := ""
	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)
	hx := `hx-target="#tbody" hx-swap="innerHTML" hx-trigger="click"`

	if r.URL.Query().Has("pagination") {
		pagination := fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/settings/app-logs?page=%d" %s>Previous</button>
        </li>`, previous, hx)

		for i := 1; i <= min(10, size); i++ {
			activeClass = config.ActiveClass(i, page)
			pagination += fmt.Sprintf(`<li class="page-item %s">
				<button class="page-link" hx-get="/settings/app-logs?page=%d" %s>%d</button>
			</li>`, activeClass, i, hx, i)
		}
		if size > 20 {
			pagination += `<li class="page-item"><button class="page-link">...</button></li>`
		}
		for i := max(11, size-11); i <= size; i++ {
			activeClass = config.ActiveClass(i, page)
			pagination += fmt.Sprintf(`<li class="page-item %s">
				<button class="page-link" hx-get="/settings/app-logs?page=%d" %s>%d</button>
			</li>`, activeClass, i, hx, i)
		}
		pagination += fmt.Sprintf(`<li class="page-item">
            <button class="page-link" hx-get="/settings/app-logs?page=%d" %s>Next</button>
        </li>`, next, hx)
		_, _ = w.Write([]byte(pagination))
		return nil
	}

	appLogs, _ := appLog.Get((current-1)*row, row, search)

	tr := ""
	for _, v := range appLogs {
		tr += fmt.Sprintf(`<tr>
            <th scope="row">%d</th>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>
				<div class="hstack gap-1">
					<button class="btn btn-danger" hx-delete="/settings/app-logs/%d"  hx-trigger='confirmed' onClick="Swal.fire({
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
        </tr>`, v.ID, v.Level, v.Message, v.CreatedAt.Format("2006-01-02 15:04:05"), v.ID)
	}
	_, _ = w.Write([]byte(tr))
	return nil
}

func (h *SettingHandler) AppLogDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	//appLog := &models.AppLog{}
	//w.Header().Set("HX-Redirect", "/setting")
	//_, _ = w.Write([]byte("success"))
	return nil
}
