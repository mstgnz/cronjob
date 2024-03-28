package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/mstgnz/cronjob/config"
)

type Auth struct{}

func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := config.Render(w, "login", map[string]any{"products": map[string]any{"test": template.HTML("<strong>test</strong>")}}, "navlink", "subscribe", "recommend", "scroll"); err != nil {
		log.Println(err)
	}
	/* if r.Method == http.MethodPost{
		// Örnek kullanıcı bilgileri (veritabanından alınır)
		username := "kullaniciAdi"
		password := "sifre123"

		// Kullanıcının gönderdiği bilgileri al
		r.ParseForm()
		submittedUsername := r.Form.Get("username")
		submittedPassword := r.Form.Get("password")

		// Kullanıcı adı ve şifre doğru mu kontrol et
		if submittedUsername != username || submittedPassword != password {
			http.Error(w, "Geçersiz kullanıcı adı veya şifre", http.StatusUnauthorized)
			return
		}

		// JWT oluşturma
		generatedToken, err := config.GenerateToken(user.ID)
		if err != nil {
			_ = config.WriteJSON(w, http.StatusOK, config.Response{Status: false, Message: "Failed to process request", Error: err.Error()})
			return
		}
		user.Token = generatedToken

		// JWT'yi tarayıcıya yazma
		cookie := &http.Cookie{
			Name:    "jwtToken",
			Value:   tokenString,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
			HttpOnly: true, // XSS saldırılarına karşı koruma için
		}
		http.SetCookie(w, cookie)

		// Kullanıcıya başarılı bir şekilde giriş yaptığını belirt
		w.Write([]byte("Başarılı giriş!"))
	} */
}

func (a *Auth) RegisterHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "register", map[string]any{"products": map[string]any{"test": template.HTML("<strong>test</strong>")}}, "navlink", "subscribe", "recommend", "scroll"); err != nil {
		log.Println(err)
	}
}

func (a *Auth) UpdateHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "profile", map[string]any{"products": map[string]any{"test": template.HTML("<strong>test</strong>")}}, "navlink", "subscribe", "recommend", "scroll"); err != nil {
		log.Println(err)
	}
}

func (a *Auth) DeleteHandler(w http.ResponseWriter, _ *http.Request) {
	if err := config.Render(w, "profile", map[string]any{"products": map[string]any{"test": template.HTML("<strong>test</strong>")}}, "navlink", "subscribe", "recommend", "scroll"); err != nil {
		log.Println(err)
	}
}
