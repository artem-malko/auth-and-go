package auth

import (
	"net/http"

	"github.com/artem-malko/auth-and-go/managers/user"
	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/api/response"

	"github.com/artem-malko/auth-and-go/api/middleware"
)

const invalidRefreshTokenErrText = "Invalid refresh token"

func (h *handlers) RefreshSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	refreshToken := middleware.GetRefreshToken(ctx)

	if refreshToken == nil {
		response.Error(w, http.StatusForbidden, invalidRefreshTokenErrText)
		return
	}

	sessionTokens, err := h.
		userManager.
		RefreshSession(*refreshToken)

	if err != nil {
		switch errors.Cause(err) {
		case user.ErrUserSessionNotFound:
			response.Error(w, http.StatusForbidden, invalidRefreshTokenErrText)
		default:
			h.logger(r).
				WithField("method", "RefreshSession").
				WithField("code", http.StatusInternalServerError).
				Error(err.Error())
			response.InternalServerError(w)

		}
		return
	}

	sessionCookies := h.createSessionCookies(*sessionTokens)
	http.SetCookie(w, &sessionCookies.accessTokenCookie)
	http.SetCookie(w, &sessionCookies.refreshTokenCookie)

	response.OKWithoutContent(w)
}
