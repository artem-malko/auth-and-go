package auth

import (
	"net/http"

	"github.com/artem-malko/auth-and-go/api/middleware"

	"github.com/artem-malko/auth-and-go/api/response"
	"github.com/artem-malko/auth-and-go/constants"
	"github.com/artem-malko/auth-and-go/infrastructure/cookie"
)

func (h *handlers) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accessTokenInfo := middleware.GetAccessTokenInfo(ctx)
	refreshToken := middleware.GetRefreshToken(ctx)

	if accessTokenInfo == nil || refreshToken == nil {
		response.Error(w, http.StatusUnauthorized, "You are not logged in")
		return
	}

	err := h.userManager.DeleteSessionBySessionID(accessTokenInfo.SessionID)

	if err != nil {
		h.logger(r).
			WithField("method", "Logout").
			WithField("code", http.StatusInternalServerError).
			Error(err.Error())
		response.InternalServerError(w)
		return
	}

	accessTokenCookie := cookie.CreateAuthTokenCookie(constants.Values.AccessTokenCookieName, h.sessionCookiesDomain)
	accessTokenCookie.MaxAge = -1
	accessTokenCookie.Value = ""
	refreshTokenCookie := cookie.CreateAuthTokenCookie(constants.Values.RefreshTokenCookieName, h.sessionCookiesDomain)
	refreshTokenCookie.MaxAge = -1
	refreshTokenCookie.Value = ""

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	response.OKWithoutContent(w)
}
