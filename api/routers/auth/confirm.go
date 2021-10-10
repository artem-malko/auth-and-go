package auth

import (
	"io"
	"net/http"

	"github.com/artem-malko/auth-and-go/managers/user"

	"github.com/artem-malko/auth-and-go/api/middleware"

	"github.com/artem-malko/auth-and-go/api/response"
	"github.com/go-chi/render"

	"github.com/google/uuid"

	"github.com/pkg/errors"
)

type confirmParams struct {
	RawConfirmationToken    string `json:"confirmation_token"`
	ParsedConfirmationToken uuid.UUID
}

func (p *confirmParams) Bind(r *http.Request) error {
	if p.RawConfirmationToken == "" {
		return errors.New("confirmation_token is required")
	}

	parsedConfirmationToken, err := uuid.Parse(p.RawConfirmationToken)

	if err != nil {
		return errors.New("confirmation_token is invalid")
	}

	p.ParsedConfirmationToken = parsedConfirmationToken

	return nil
}

func (h *handlers) ConfirmRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Requests from trusted clients can be proceeded only
	if !middleware.IsTrustedClient(ctx) {
		response.Error(w, http.StatusForbidden, "Unknown client")
		return
	}

	reqParams := &confirmParams{}

	if err := render.Bind(r, reqParams); err != nil {
		switch errors.Cause(err) {
		case io.EOF:
			response.Error(w, http.StatusBadRequest, "Request body is required. Pass correct confirmation_token.")
		default:
			response.Error(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	u, sessionTokens, err := h.userManager.ConfirmRegistration(reqParams.ParsedConfirmationToken)

	if err != nil {
		switch errors.Cause(err) {
		case user.ErrUserIncorrectTokenToUse:
			response.Error(w, http.StatusForbidden, "Token can not be used")
		default:
			h.logger(r).
				WithField("method", "ConfirmRegistration").
				WithField("code", http.StatusInternalServerError).
				Error(err.Error())
			response.InternalServerError(w)
		}
		return
	}

	sessionCookies := h.createSessionCookies(*sessionTokens)
	http.SetCookie(w, &sessionCookies.accessTokenCookie)
	http.SetCookie(w, &sessionCookies.refreshTokenCookie)

	response.OK(w, u)
}
