package handler

import (
	"log"
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type Home struct{}

func (h *Home) HomeHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "home", map[string]any{}); err != nil {
		log.Println(err)
	}
}

func PostHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "post", map[string]any{}); err != nil {
		log.Println(err)
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := config.Render(w, "create", map[string]any{}); err != nil {
			log.Println(err)
		}
	} else {
		title := r.Form.Get("title")
		fullname := r.Form.Get("fullname")
		content := r.Form.Get("content")

		if title != "" && fullname != "" && content != "" {

			w.Write([]byte("success"))
		} else {
			w.Write([]byte("invalid form"))
		}
	}

}
