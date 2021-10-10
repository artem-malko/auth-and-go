package users

import (
	"html"
	"io"
	"net/http"
	"net/mail"
	"strings"

	"github.com/artem-malko/auth-and-go/models"

	"github.com/artem-malko/auth-and-go/api/response"
	userManager "github.com/artem-malko/auth-and-go/managers/user"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

var addressParser = &mail.AddressParser{}

type CreateUserParams struct {
	Email          string          `json:"email"`
	Password       string          `json:"password"`
	PasswordRepeat string          `json:"password_repeat"`
	RawClientID    string          `json:"client_id"`
	ClientID       models.ClientID `json:"-"`
}

func (c *CreateUserParams) Bind(_ *http.Request) error {
	c.Email = html.EscapeString(strings.TrimSpace(strings.ToLower(c.Email)))

	if c.Email == "" {
		return errors.New("Email is required")
	}

	_, err := addressParser.Parse(c.Email)

	if err != nil {
		return errors.New("Valid email is required *@*.*")
	}

	if c.Password == "" {
		return errors.New("Password is required")
	}

	if c.Password != c.PasswordRepeat {
		return errors.New("Password and repeated password are not equal")
	}

	if c.RawClientID == "" {
		return errors.New("client_id is required")
	}

	if !models.CheckClientID(c.RawClientID) {
		return errors.New("incorrect client_id")
	}

	c.ClientID = models.ClientID(c.RawClientID)

	return nil
}

func (h *handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	reqParams := &CreateUserParams{}

	if err := render.Bind(r, reqParams); err != nil {
		switch errors.Cause(err) {
		case io.EOF:
			response.Error(w, http.StatusBadRequest, "Request body is required. Pass correct email and password.")
		default:
			response.Error(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	err := h.userManager.CreateUserWithEmail(reqParams.Email, reqParams.Password, reqParams.ClientID)

	if err != nil {
		switch errors.Cause(err) {
		case userManager.ErrUserWithSameIdentityExists:
			// Its ok, cause its much more safe to send ok on each creation for existed account
			response.OKWithoutContent(w)
		default:
			h.logger(r).
				WithField("method", "CreateUser").
				WithField("code", http.StatusInternalServerError).
				Error(err.Error())
			response.InternalServerError(w)
		}
		return
	}

	response.OKWithoutContent(w)
}
