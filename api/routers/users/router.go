package users

import (
	"net/http"

	"github.com/artem-malko/auth-and-go/api/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/apex/log"
	"github.com/artem-malko/auth-and-go/api/response"
	"github.com/artem-malko/auth-and-go/managers/user"
)

type RouterConfig struct {
	Logger               log.Interface
	UserManager          user.Manager
	SessionCookiesDomain string
	AccessTokenSecretKey string
}

// handlers include all handlers for API
type handlers struct {
	logger               func(r *http.Request) log.Interface
	userManager          user.Manager
	sessionCookiesDomain string
	accessTokenSecretKey string
}

func NewRouter(cfg RouterConfig) func(userRouter chi.Router) {
	handlers := &handlers{
		logger:               response.CreateBoundLogger(cfg.Logger),
		userManager:          cfg.UserManager,
		sessionCookiesDomain: cfg.SessionCookiesDomain,
		accessTokenSecretKey: cfg.AccessTokenSecretKey,
	}

	accessTokenMiddleware := middleware.CreateAccessTokenMiddleware(cfg.UserManager, cfg.AccessTokenSecretKey)

	return func(userRouter chi.Router) {
		const userIDContextKey = "userID"
		const userIDPathParamName = "id"
		const userNameQueryParamName = "username"

		correctUserID := middleware.CreateCorrectUUIDInParamMiddleware(
			cfg.Logger,
			userIDContextKey,
			userIDPathParamName,
		)

		getUserByID := handlers.CreateGetUserByID(userIDContextKey)
		getUserByName := handlers.CreateGetUserByName(userNameQueryParamName)

		// All users
		userRouter.With(correctUserID).Get("/{"+userIDPathParamName+"}", getUserByID)
		userRouter.Get("/", getUserByName)

		userRouter.With(accessTokenMiddleware).Route("/me", func(userMeRouter chi.Router) {
			// Own user
			// It does not need accessTokenMiddleware, but its ok to be here
			userMeRouter.Post("/", handlers.CreateUser)

			userMeRouter.Get("/", handlers.GetFullUser)

			userMeRouter.Get("/accepted_activities", handlers.Stub)
			userMeRouter.Get("/analytics", handlers.Stub)

			userMeRouter.Delete("/", handlers.DeleteUser)

			userMeRouter.Patch("/profile", handlers.Stub)
			userMeRouter.Patch("/username", handlers.UpdateUserName)
		})
	}
}
