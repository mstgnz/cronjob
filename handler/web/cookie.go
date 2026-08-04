package web

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mstgnz/cronjob/pkg/auth"
)

const authCookieName = "Authorization"

// setAuthCookie stores the session token.
// HttpOnly keeps the token out of reach of any script running on the page, SameSite
// blocks the cookie from riding along on cross site state changing requests, and
// Secure is enabled outside development so it is never sent over plain HTTP.
func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    strings.Join([]string{"Bearer", token}, " "),
		Path:     "/",
		Expires:  time.Now().Add(auth.TokenTTL),
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	})
}

// clearAuthCookie expires the session cookie. The attributes have to match the ones
// used when setting it, otherwise the browser keeps the original cookie in place.
func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	})
}
