package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/service"
)

type Web struct {
	service.Api
}

func (wb *Web) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := config.Render(w, "login", map[string]any{}); err != nil {
			log.Println(err)
		}
	case http.MethodPost:
		_ = config.WriteJSON(w, http.StatusAccepted, config.Response{Status: true, Message: "login", Data: r.Body})
	}
}

func (wb *Web) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{}
	if r.Method == http.MethodPost {
		_, response := wb.Api.RegisterService(w, r)
		if response.Status {
			token, _ := response.Data.(string)
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(24 * time.Hour),
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		data["status"] = response.Status
		data["message"] = response.Message
	}
	if err := config.Render(w, "register", data); err != nil {
		log.Println("RegisterHandler Error: ", err)
	}
}

func (wb *Web) HomeHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "home", map[string]any{}); err != nil {
		log.Println(err)
	}
}

func (wb *Web) ListHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "schedule", map[string]any{}); err != nil {
		log.Println(err)
	}
}

func (wb *Web) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if err := config.Render(w, "profile", map[string]any{}); err != nil {
		log.Println(err)
	}
}
