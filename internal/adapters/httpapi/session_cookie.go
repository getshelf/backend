package httpapi

import (
	"net/http"
	"time"
)

const sessionCookieName = "shelf_session"

func newSessionCookie(token string, expiresAt time.Time, secure bool) http.Cookie {
	return http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt.UTC(),
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func clearSessionCookie(secure bool) http.Cookie {
	cookie := newSessionCookie("", time.Unix(0, 0), secure)
	cookie.MaxAge = -1
	return cookie
}
