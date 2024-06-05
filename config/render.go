package config

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func Render(w http.ResponseWriter, page string, data map[string]any, partials ...string) error {

	var t *template.Template
	var err error

	if len(partials) > 0 {
		partialPaths := make([]string, len(partials))
		for i, partial := range partials {
			partialPaths[i] = path.Join("./viwes/components", fmt.Sprintf("%s.gohtml", partial))
		}

		templateFiles := append(partialPaths, path.Join("./viwes/pages", fmt.Sprintf("%s.gohtml", page)), path.Join("./viwes/index.gohtml"))
		t, err = template.ParseFiles(templateFiles...)
	} else {
		t, err = template.ParseFiles(path.Join("./viwes/pages", fmt.Sprintf("%s.gohtml", page)), path.Join("./viwes/index.gohtml"))
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
