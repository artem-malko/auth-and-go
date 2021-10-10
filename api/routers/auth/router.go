package auth

import (
	"net/http"

	"golang.org/x/oauth2/facebook"

	"github.com/artem-malko/auth-and-go/api/middleware"

	"golang.org/x/oauth2/google"

	"golang.org/x/oauth2"

	"github.com/apex/log"
	"github.com/artem-malko/auth-and-go/api/response"
	"github.com/artem-malko/auth-and-go/managers/user"
	"github.com/go-chi/chi/v5"
)

type RouterConfig struct {
	Logger                          log.Interface
	UserManager                     user.Manager
	SessionCookiesDomain            string
	AccessTokenSecretKey            string
	RefreshTokenSecretKey           string
	GoogleAuthClientID              string
	GoogleAuthClientSecret          string
	GoogleAuthCallbackRedirectURL   string
	FacebookAuthClientID            string
	FacebookAuthClientSecret        string
	FacebookAuthCallbackRedirectURL string
	FrontendAppURL                  string
}

// handlers include all handlers for API
type handlers struct {
	logger                func(r *http.Request) log.Interface
	userManager           user.Manager
	sessionCookiesDomain  string
	accessTokenSecretKey  string
	refreshTokenSecretKey string
	googleOAuthConfig     oauth2.Config
	facebookOAuthConfig   oauth2.Config
	frontendAppURL        string
}

func NewRouter(cfg RouterConfig) func(userRouter chi.Router) {
	handlers := &handlers{
		logger:               response.CreateBoundLogger(cfg.Logger),
		userManager:          cfg.UserManager,
		sessionCookiesDomain: cfg.SessionCookiesDomain,
		googleOAuthConfig: oauth2.Config{
			RedirectURL:  cfg.GoogleAuthCallbackRedirectURL,
			ClientID:     cfg.GoogleAuthClientID,
			ClientSecret: cfg.GoogleAuthClientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
		facebookOAuthConfig: oauth2.Config{
			RedirectURL:  cfg.FacebookAuthCallbackRedirectURL,
			ClientID:     cfg.FacebookAuthClientID,
			ClientSecret: cfg.FacebookAuthClientSecret,
			Scopes: []string{
				"public_profile",
				"email",
			},
			Endpoint: facebook.Endpoint,
		},
		accessTokenSecretKey:  cfg.AccessTokenSecretKey,
		refreshTokenSecretKey: cfg.RefreshTokenSecretKey,
		frontendAppURL:        cfg.FrontendAppURL,
	}

	accessTokenMiddleware := middleware.CreateAccessTokenMiddleware(cfg.UserManager, cfg.AccessTokenSecretKey)
	refreshTokenMiddleware := middleware.CreateRefreshTokenMiddleware(cfg.RefreshTokenSecretKey)

	return func(authRouter chi.Router) {
		authRouter.Post("/login", handlers.LoginByEmailAndPassword)
		authRouter.With(accessTokenMiddleware, refreshTokenMiddleware).Post("/logout", handlers.Logout)
		authRouter.With(accessTokenMiddleware, refreshTokenMiddleware).Post("/refresh", handlers.RefreshSession)
		authRouter.Post("/confirm", handlers.ConfirmRegistration)
		authRouter.Get("/autologin", handlers.Stub)

		authRouter.Get("/google/login", handlers.CreateOAuthLogin(handlers.googleOAuthConfig))
		authRouter.Get("/google/callback", handlers.CreateOAuthCallback(handlers.googleOAuthConfig, "google"))

		authRouter.Get("/facebook/login", handlers.CreateOAuthLogin(handlers.facebookOAuthConfig))
		authRouter.Get("/facebook/callback", handlers.CreateOAuthCallback(handlers.facebookOAuthConfig, "facebook"))
	}
}
