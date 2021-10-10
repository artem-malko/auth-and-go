package users

import (
	"io"
	"net/http"

	"github.com/artem-malko/auth-and-go/api/middleware"

	"github.com/artem-malko/auth-and-go/managers/user"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/api/response"
	"github.com/go-chi/render"
)

type UpdateUserNameParams struct {
	UserName string `json:"user_name"`
}

func (c *UpdateUserNameParams) Bind(_ *http.Request) error {
	if c.UserName == "" {
		return errors.New("user_name is required")
	}

	return nil
}
func (h *handlers) UpdateUserName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqParams := &UpdateUserNameParams{}

	if err := render.Bind(r, reqParams); err != nil {
		switch errors.Cause(err) {
		case io.EOF:
			response.Error(w, http.StatusBadRequest, "user_name is required")
		default:
			response.Error(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	accessTokenInfo := middleware.GetAccessTokenInfo(ctx)

	if accessTokenInfo == nil {
		response.Error(w, http.StatusUnauthorized, "You are not authed")
		return
	}

	err := h.userManager.UpdateAccountName(accessTokenInfo.AccountID, reqParams.UserName)

	if err != nil {
		switch errors.Cause(err) {
		case user.ErrUserIsNotUpdated:
			response.Error(w, http.StatusUnprocessableEntity, "There is no user with ID "+accessTokenInfo.AccountID.String())
		case user.ErrUserWithSameNameExists:
			response.Error(w, http.StatusConflict, "user_name '"+reqParams.UserName+"' already exists. Choose another name.")
		default:
			h.logger(r).
				WithField("method", "UpdateUserName").
				WithField("code", http.StatusInternalServerError).
				Error(err.Error())
			response.InternalServerError(w)
		}
		return
	}

	response.OKWithoutContent(w)
}
