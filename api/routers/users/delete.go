package users

import (
	"net/http"

	"github.com/artem-malko/auth-and-go/constants"

	"github.com/artem-malko/auth-and-go/infrastructure/cookie"

	"github.com/artem-malko/auth-and-go/api/middleware"

	"github.com/artem-malko/auth-and-go/api/response"
	userManager "github.com/artem-malko/auth-and-go/managers/user"
	"github.com/pkg/errors"
)

func (h *handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessTokenInfo := middleware.GetAccessTokenInfo(ctx)

	if accessTokenInfo == nil {
		response.Error(w, http.StatusUnauthorized, "You are not authed")
		return
	}

	err := h.userManager.DeleteUserByID(accessTokenInfo.AccountID)

	if err != nil {
		switch errors.Cause(err) {
		case userManager.ErrUserIsNotUpdated:
			response.Error(
				w,
				http.StatusUnprocessableEntity,
				"There is no user with ID "+accessTokenInfo.AccountID.String(),
			)
		default:
			h.logger(r).
				WithField("method", "DeleteUser").
				WithField("code", http.StatusInternalServerError).
				Error(err.Error())
			response.InternalServerError(w)
		}
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
