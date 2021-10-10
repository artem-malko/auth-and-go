package api

import (
	"net/http"

	"github.com/artem-malko/auth-and-go/api/routers/auth"

	"github.com/artem-malko/auth-and-go/managers/user"

	"github.com/apex/log"
	"github.com/artem-malko/auth-and-go/api/middleware"
	"github.com/artem-malko/auth-and-go/api/response"
	"github.com/artem-malko/auth-and-go/api/routers/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Config struct {
	AllowedOrigins                  string
	AccessTokenSecretKey            string
	RefreshTokenSecretKey           string
	SessionCookiesDomain            string
	GoogleAuthClientID              string
	GoogleAuthClientSecret          string
	GoogleAuthCallbackRedirectURL   string
	FacebookAuthClientID            string
	FacebookAuthClientSecret        string
	FacebookAuthCallbackRedirectURL string
	NeedToCheckClient               bool
	FrontendAppURL                  string
}

// API struct for API instance
type API struct {
	logger       log.Interface
	userManager  user.Manager
	healthRouter chi.Router
	config       Config
}

// NewRouter create API
func New(
	logger log.Interface,
	userManager user.Manager,
	healthRouter chi.Router,
	config Config,
) API {
	return API{
		logger:       logger,
		healthRouter: healthRouter,
		userManager:  userManager,
		config:       config,
	}
}

// CreateRouter init router for App
func (a API) CreateRouter() http.Handler {
	router := chi.NewRouter()

	// debug.GenerateData(a.userManager, a.challengesManager, a.activityManager)

	usersRouter := users.NewRouter(
		users.RouterConfig{
			Logger:               a.logger,
			UserManager:          a.userManager,
			SessionCookiesDomain: a.config.SessionCookiesDomain,
			AccessTokenSecretKey: a.config.AccessTokenSecretKey,
		},
	)
	authRouter := auth.NewRouter(
		auth.RouterConfig{
			Logger:                          a.logger,
			UserManager:                     a.userManager,
			SessionCookiesDomain:            a.config.SessionCookiesDomain,
			AccessTokenSecretKey:            a.config.AccessTokenSecretKey,
			RefreshTokenSecretKey:           a.config.RefreshTokenSecretKey,
			GoogleAuthClientID:              a.config.GoogleAuthClientID,
			GoogleAuthClientSecret:          a.config.GoogleAuthClientSecret,
			GoogleAuthCallbackRedirectURL:   a.config.GoogleAuthCallbackRedirectURL,
			FacebookAuthClientID:            a.config.FacebookAuthClientID,
			FacebookAuthClientSecret:        a.config.FacebookAuthClientSecret,
			FacebookAuthCallbackRedirectURL: a.config.FacebookAuthCallbackRedirectURL,
			FrontendAppURL:                  a.config.FrontendAppURL,
		},
	)

	corsInstance := cors.New(cors.Options{
		AllowedOrigins: []string{a.config.AllowedOrigins},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders: []string{
			"Accept",
			"Keep-Alive",
			"Cache-Control",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           3600,
	})

	router.Use(corsInstance.Handler)

	// Top level router for all version of API
	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(
			middleware.CreateClientIPMiddleware(a.logger),
			middleware.RequestID,
			middleware.CreateRecoverMiddleware(a.logger),
			middleware.CreateAccessLogMiddleware(a.logger),
			middleware.CreateCheckClientMiddleware(a.config.NeedToCheckClient),
		)

		apiRouter.Mount("/healthcheck", a.healthRouter)

		// v1.0 router
		apiRouter.Route("/v1.0", func(v1Router chi.Router) {

			// User
			v1Router.Route("/users", usersRouter)

			// Auth
			v1Router.Route("/auth", authRouter)
		})

		apiRouter.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			response.NotAllowed(w)
		})
	})

	// Default handlers
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.NotFound(w)
	})

	return router
}
