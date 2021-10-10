package auth

import (
	"html"
	"io"
	"net/http"
	"net/mail"
	"strings"

	"github.com/artem-malko/auth-and-go/models"

	"github.com/artem-malko/auth-and-go/api/middleware"

	"github.com/artem-malko/auth-and-go/managers/user"

	"github.com/pkg/errors"

	"github.com/go-chi/render"

	"github.com/artem-malko/auth-and-go/api/response"
)

var addressParser = &mail.AddressParser{}

type loginByEmailAndPasswordParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
}

func (c *loginByEmailAndPasswordParams) Bind(r *http.Request) error {
	c.Email = html.EscapeString(strings.TrimSpace(strings.ToLower(c.Email)))

	if c.Email == "" {
		return errors.New("Email is required")
	}

	_, err := addressParser.Parse(c.Email)

	if err != nil {
		return errors.New("Valid email is required *@*.*")
	}

	if c.Password == "" {
		return errors.New("password is required")
	}

	if c.ClientID == "" {
		return errors.New("client_id is required")
	}

	if !models.CheckClientID(c.ClientID) {
		return errors.New("incorrect client_id")
	}

	return nil
}

// Add login via username and password
func (h *handlers) LoginByEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	reqParams := &loginByEmailAndPasswordParams{}

	if err := render.Bind(r, reqParams); err != nil {
		switch errors.Cause(err) {
		case io.EOF:
			response.Error(w, http.StatusBadRequest, "Request body is required. Pass correct email and password.")
		default:
			response.Error(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	ctx := r.Context()
	clientIP := middleware.GetClientIP(ctx)
	clientID := models.ClientID(reqParams.ClientID)

	u, sessionTokens, err := h.
		userManager.
		LoginWithEmailAndPassword(reqParams.Email, reqParams.Password, clientIP, clientID)

	if err != nil {
		switch errors.Cause(err) {
		case user.ErrUserNoIdentitiesFound:
			response.Error(w, http.StatusForbidden, "Incorrect email or password")
		default:
			h.logger(r).
				WithField("method", "LoginByEmailAndPassword").
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
