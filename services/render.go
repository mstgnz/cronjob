package services

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/models"
)

func Render(w http.ResponseWriter, r *http.Request, page string, data map[string]any, partials ...string) error {

	var t *template.Template
	var err error

	cUser, ok := r.Context().Value(config.CKey("user")).(*models.User)

	if data == nil {
		data = make(map[string]any)
	}
	if cUser != nil {
		data["auth"] = cUser
		data["isAuth"] = ok
	}

	if len(partials) > 0 {
		partialPaths := make([]string, len(partials))
		for i, partial := range partials {
			partialPaths[i] = path.Join("./views/components", fmt.Sprintf("%s.gohtml", partial))
		}

		templateFiles := append(partialPaths, path.Join("./views/pages", fmt.Sprintf("%s.gohtml", page)), path.Join("./views/index.gohtml"))
		t, err = template.ParseFiles(templateFiles...)
	} else {
		t, err = template.ParseFiles(path.Join("./views/pages", fmt.Sprintf("%s.gohtml", page)), path.Join("./views/index.gohtml"))
	}

	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return err
	}

	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
