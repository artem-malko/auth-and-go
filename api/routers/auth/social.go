package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/artem-malko/auth-and-go/api/middleware"

	"github.com/artem-malko/auth-and-go/managers/user"
	"github.com/artem-malko/auth-and-go/models"

	"github.com/pkg/errors"

	"golang.org/x/oauth2"

	"github.com/artem-malko/auth-and-go/api/response"
)

const clientIDCookieName = "_oauth_client_id"
const successRedirectCookieName = "_success_login_redirect_url"
const unsuccessfulRedirectCookieName = "_unsuccessful_login_redirect_url"
const oAuthCookieName = "_oauth_state"

type socialData struct {
	ID              string
	Email           string
	IsEmailVerified bool
	FirstName       string
	LastName        string
	AvatarURL       string
}

func (h *handlers) CreateOAuthLogin(oAuthConfig oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create oauthState cookie
		oauthState := generateStateOauthCookie(w)
		expiration := time.Now().Add(20 * time.Minute)
		clientID := "web"

		if clientIDFromQuery := r.URL.Query().Get("client_id"); clientIDFromQuery != "" {
			clientID = clientIDFromQuery
		}

		clientIDCookie := http.Cookie{
			Name:    clientIDCookieName,
			Value:   clientID,
			Expires: expiration,
		}
		http.SetCookie(w, &clientIDCookie)

		if redirectURL := r.URL.Query().Get("success_login_redirect_url"); redirectURL != "" {
			redirectURLCookie := http.Cookie{
				Name:    successRedirectCookieName,
				Value:   redirectURL,
				Expires: expiration,
			}
			http.SetCookie(w, &redirectURLCookie)
		}

		if redirectURL := r.URL.Query().Get("unsuccessful_login_redirect_url"); redirectURL != "" {
			redirectURLCookie := http.Cookie{
				Name:    unsuccessfulRedirectCookieName,
				Value:   redirectURL,
				Expires: expiration,
			}
			http.SetCookie(w, &redirectURLCookie)
		}

		/*
			AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
			You must always provide a non-empty string and
			validate that it matches the state query parameter on your redirect callback.
		*/
		u := oAuthConfig.AuthCodeURL(oauthState)

		response.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

func (h *handlers) CreateOAuthCallback(oAuthConfig oauth2.Config, networkType models.SocialNetworkType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientIP := middleware.GetClientIP(ctx)
		params := h.getSocialLoginParams(w, r)
		// Read oauthState from Cookie
		oauthState, err := r.Cookie(oAuthCookieName)

		if err == http.ErrNoCookie || oauthState == nil || r.FormValue("state") != oauthState.Value {
			h.logger(r).
				WithField("method", "OAuthCallback").
				WithField("social_network", networkType).
				Error(errors.Wrap(err, "Invalid oauth state").Error())
			http.Redirect(w, r, params.UnsuccessfulRedirectURL, http.StatusTemporaryRedirect)
			return
		}

		data, err := h.getUserData(r.FormValue("code"), oAuthConfig, getSocialNetworkAPIURL(networkType))

		if err != nil {
			h.logger(r).
				WithField("method", "OAuthCallback").
				WithField("social_network", networkType).
				Error(errors.Wrap(err, "Error during getUserData").Error())

			http.Redirect(w, r, params.UnsuccessfulRedirectURL, http.StatusTemporaryRedirect)
			return
		}

		parsedData, err := parseUserData(networkType, data)

		if err != nil {
			h.logger(r).
				WithField("method", "OAuthCallback").
				WithField("social_network", networkType).
				Error(errors.Wrap(err, "Error during unmarshal social data").Error())

			http.Redirect(w, r, params.UnsuccessfulRedirectURL, http.StatusTemporaryRedirect)
			return
		}

		sessionTokens, err := h.userManager.ContinueWithOAuth(
			user.ContinueWithOAuthParams{
				SocialNetworkType: networkType,
				SocialID:          parsedData.ID,
				Email:             parsedData.Email,
				IsEmailVerified:   parsedData.IsEmailVerified,
				FirstName:         parsedData.FirstName,
				LastName:          parsedData.LastName,
				AvatarURL:         parsedData.AvatarURL,
				ClientID:          params.ClientID,
				ClientIP:          clientIP,
			},
		)

		if err != nil {
			h.logger(r).
				WithField("method", "OAuthCallback").
				WithField("social_network", networkType).
				Error(errors.Wrap(err, "Error during LoginWithOAuth").Error())

			http.Redirect(w, r, params.UnsuccessfulRedirectURL, http.StatusTemporaryRedirect)
			return
		}

		sessionCookies := h.createSessionCookies(*sessionTokens)

		http.SetCookie(w, &sessionCookies.accessTokenCookie)
		http.SetCookie(w, &sessionCookies.refreshTokenCookie)

		response.Redirect(w, r, params.SuccessRedirectURL, http.StatusTemporaryRedirect)
	}
}

type socialLoginParams struct {
	SuccessRedirectURL      string
	UnsuccessfulRedirectURL string
	ClientID                models.ClientID
}

func (h *handlers) getSocialLoginParams(w http.ResponseWriter, r *http.Request) socialLoginParams {
	params := socialLoginParams{
		SuccessRedirectURL:      h.frontendAppURL,
		UnsuccessfulRedirectURL: h.frontendAppURL,
		ClientID:                "web",
	}

	successLoginRedirectCookie, err := r.Cookie(successRedirectCookieName)

	if err == nil && successLoginRedirectCookie != nil && successLoginRedirectCookie.Value != "" {
		params.SuccessRedirectURL = successLoginRedirectCookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   successRedirectCookieName,
			Value:  "",
			MaxAge: -1,
		})
	}

	unsuccessfulLoginRedirectCookie, err := r.Cookie(unsuccessfulRedirectCookieName)

	if err == nil && unsuccessfulLoginRedirectCookie != nil && unsuccessfulLoginRedirectCookie.Value != "" {
		params.UnsuccessfulRedirectURL = unsuccessfulLoginRedirectCookie.Value
		http.SetCookie(w, &http.Cookie{
			Name:   unsuccessfulRedirectCookieName,
			Value:  "",
			MaxAge: -1,
		})
	}

	clientIDCookie, err := r.Cookie(clientIDCookieName)

	if err == nil && clientIDCookie != nil && clientIDCookie.Value == "" && !models.CheckClientID(clientIDCookie.Value) {
		params.ClientID = models.ClientID(clientIDCookie.Value)
		http.SetCookie(w, &http.Cookie{
			Name:   clientIDCookieName,
			Value:  "",
			MaxAge: -1,
		})
	}

	return params
}

/*
	Exchange code to access_token and request user data
*/
func (h *handlers) getUserData(code string, oAuthConfig oauth2.Config, url string) ([]byte, error) {
	// Use code to get token and get user info
	token, err := oAuthConfig.Exchange(context.Background(), code)

	if err != nil {
		return nil, errors.Wrap(err, "code exchange is wrong")
	}

	resp, err := http.Get(url + token.AccessToken)

	if err != nil {
		return nil, errors.Wrap(err, "failed getting user info")
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "failed read resp")
	}

	return contents, nil
}

func parseUserData(networkType models.SocialNetworkType, rawData []byte) (*socialData, error) {
	switch networkType {
	case "facebook":
		return parseUserDataFromFacebook(rawData)
	case "google":
		fallthrough
	default:
		return parseUserDataFromGoogle(rawData)
	}
}

type googleData struct {
	ID              string `json:"sub"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"verified_email"`
	FirstName       string `json:"given_name"`
	LastName        string `json:"family_name"`
	AvatarURL       string `json:"picture"`
}

func parseUserDataFromGoogle(rawData []byte) (*socialData, error) {
	googleData := new(googleData)

	err := json.Unmarshal(rawData, &googleData)

	if err != nil {
		return nil, err
	}

	return &socialData{
		ID:              googleData.ID,
		Email:           googleData.Email,
		IsEmailVerified: true,
		FirstName:       googleData.FirstName,
		LastName:        googleData.LastName,
		AvatarURL:       googleData.AvatarURL,
	}, nil
}

type facebookData struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func parseUserDataFromFacebook(rawData []byte) (*socialData, error) {
	facebookData := new(facebookData)

	err := json.Unmarshal(rawData, &facebookData)

	if err != nil {
		return nil, err
	}

	return &socialData{
		ID:              facebookData.ID,
		Email:           facebookData.Email,
		IsEmailVerified: true,
		FirstName:       facebookData.FirstName,
		LastName:        facebookData.LastName,
		AvatarURL:       "https://graph.facebook.com/v6.0/" + facebookData.ID + "/picture?width=1024",
	}, nil
}

func getSocialNetworkAPIURL(networkType models.SocialNetworkType) string {
	switch networkType {
	case "facebook":
		return "https://graph.facebook.com/v6.0/me?fields=first_name,last_name,email,gender,name&access_token="
	case "google":
		fallthrough
	default:
		return "https://www.googleapis.com/oauth2/v3/userinfo?access_token="
	}
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(20 * time.Minute)
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: oAuthCookieName, Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}
