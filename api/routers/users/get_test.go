package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/artem-malko/auth-and-go/api/response"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/credentials"

	"github.com/artem-malko/auth-and-go/constants"
	"github.com/artem-malko/auth-and-go/infrastructure/cookie"

	"github.com/artem-malko/auth-and-go/infrastructure/logger"

	"github.com/go-chi/chi/v5"

	"github.com/artem-malko/auth-and-go/managers/user"

	. "github.com/artem-malko/auth-and-go/forks/goblin"
)

func TestGetFullUser(t *testing.T) {
	g := Goblin(t)

	var server *httptest.Server
	var mockedUserManager *user.MockManager
	var httpClient *http.Client
	accessTokenSecretKey := "dwmanjmwofien8rctcjnvahvluq7ch2r"
	sessionCookiesDomain := "."

	g.Describe("users/me GetFullUser", func() {
		g.BeforeEach(func() {
			router := chi.NewRouter()
			httpClient = &http.Client{}
			mockedUserManager = new(user.MockManager)
			userRouter := NewRouter(RouterConfig{
				Logger:               logger.NewForTests(),
				UserManager:          mockedUserManager,
				SessionCookiesDomain: sessionCookiesDomain,
				AccessTokenSecretKey: accessTokenSecretKey,
			})
			router.Route("/users", userRouter)
			server = httptest.NewServer(router)
		})

		g.It("Get users/me without accessToken return StatusUnauthorized code", func() {
			resp, err := httpClient.Get(server.URL + "/users/me")

			if err != nil {
				g.Fail(err)
			}

			data, _ := ioutil.ReadAll(resp.Body)
			var parsedData response.ErrorResponse
			_ = json.Unmarshal(data, &parsedData)

			g.Assert(resp.StatusCode).Equal(http.StatusUnauthorized)
			g.Assert(parsedData.ErrorData.Code).Equal(http.StatusUnauthorized)
			g.Assert(parsedData.ErrorData.Message).Equal("You are not authed")
			g.Assert(parsedData.ErrorData.Data).Equal(nil)
		})

		g.It("Get users/me with valid accessToken return StatusOK and user data", func() {
			accountID := uuid.New()
			accessToken := uuid.New()
			session := models.Session{
				AccountID:              accountID,
				AccessToken:            accessToken,
				AccessTokenExpiresDate: time.Now().Add(20 * time.Minute),
			}
			mockedUserManager.On("GetSessionByAccessToken", accessToken).Return(&session, nil).Maybe()
			mockedUserManager.On("GetFullUser", accountID).Return(&models.User{ID: accountID}, nil)

			req, _ := http.NewRequest("GET", server.URL+"/users/me", nil)
			accessTokenString, _ := credentials.CreateAccessToken(session, []byte(accessTokenSecretKey))
			accessTokenCookie := cookie.CreateAuthTokenCookie(
				constants.Values.AccessTokenCookieName,
				sessionCookiesDomain,
			)
			accessTokenCookie.MaxAge = constants.Values.AccessTokenMaxAgeInSeconds
			accessTokenCookie.Value = accessTokenString
			req.AddCookie(&accessTokenCookie)

			resp, err := httpClient.Do(req)

			if err != nil {
				g.Fail(err)
			}

			data, _ := ioutil.ReadAll(resp.Body)
			var parsedData response.SuccessResponse
			_ = json.Unmarshal(data, &parsedData)
			var userID uuid.UUID

			for k, v := range parsedData.Data.(map[string]interface{}) {
				if k == "id" {
					userID = uuid.MustParse(v.(string))
				}
			}

			g.Assert(resp.StatusCode).Equal(http.StatusOK)
			g.Assert(userID).Equal(accountID)

			mockedUserManager.AssertExpectations(t)
		})
	})
}
