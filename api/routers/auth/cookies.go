package auth

import (
	"net/http"

	"github.com/artem-malko/auth-and-go/constants"
	"github.com/artem-malko/auth-and-go/infrastructure/cookie"
	"github.com/artem-malko/auth-and-go/models"
)

type sessionCookies struct {
	accessTokenCookie  http.Cookie
	refreshTokenCookie http.Cookie
}

func (h *handlers) createSessionCookies(sessionTokens models.SessionTokens) sessionCookies {
	accessTokenCookie := cookie.CreateAuthTokenCookie(constants.Values.AccessTokenCookieName, h.sessionCookiesDomain)
	accessTokenCookie.MaxAge = constants.Values.AccessTokenMaxAgeInSeconds
	accessTokenCookie.Value = sessionTokens.AccessToken
	refreshTokenCookie := cookie.CreateAuthTokenCookie(constants.Values.RefreshTokenCookieName, h.sessionCookiesDomain)
	refreshTokenCookie.MaxAge = constants.Values.RefreshTokenMaxAgeInSeconds
	refreshTokenCookie.Value = sessionTokens.RefreshToken

	return sessionCookies{
		accessTokenCookie,
		refreshTokenCookie,
	}
}
