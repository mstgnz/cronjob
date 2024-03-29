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
