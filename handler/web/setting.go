package web

import (
	"math"
	"net/http"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/services"
)

type SettingHandler struct {
	user *services.UserService
}

func (h *SettingHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return services.Render(w, r, "setting", map[string]any{})
}

func (h *SettingHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) error {
	code, response := h.user.RegisterService(w, r)
	if response.Status && code == http.StatusCreated {
		w.Header().Set("HX-Redirect", "/setting")
	}
	_, _ = w.Write([]byte(response.Message))
	return nil
}

func (h *SettingHandler) UserHandler(w http.ResponseWriter, r *http.Request) error {
	user := &models.User{}

	total := user.Count()
	row := 1

	page := config.GetIntQuery(r, "page")
	size := int(math.Ceil(float64(total) / float64(row)))

	current := config.Clamp(page, 1, size)
	previous := config.Clamp(current-1, 1, size)
	next := config.Clamp(current+1, 1, size)

	users := user.Users((current-1)*row, row)

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
