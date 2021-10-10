package cookie

import "net/http"

func CreateAuthTokenCookie(name string, domain string) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		// 1 session by default
		MaxAge:   0,
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
	}
}
