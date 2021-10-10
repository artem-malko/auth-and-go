package config

import "gopkg.in/alecthomas/kingpin.v2"

// @TODO_ARTEM

// Config type for config of whole application
type Config struct {
	Log struct {
		Level  string
		Format string
	}

	Server struct {
		Port string
	}

	Postgres struct {
		DSN string
	}

	Cors struct {
		AllowedOrigins string
	}

	Auth struct {
		SessionCookiesDomain            string
		AccessTokenSecretKey            string
		RefreshTokenSecretKey           string
		GoogleAuthClientID              string
		GoogleAuthClientSecret          string
		GoogleAuthCallbackRedirectURL   string
		FacebookAuthClientID            string
		FacebookAuthClientSecret        string
		FacebookAuthCallbackRedirectURL string
	}

	FrontendAppURL string

	NeedToCheckClient bool

	S3 struct {
		AccessKeyID     string
		SecretAccessKey string
		AWSRegion       string
		BucketName      string
	}

	MediaFilesPrefix string
}

// New create config which is based cli flags and env vars
func New() Config {
	cfg := Config{}

	kingpin.Flag("log-level", "Log level").
		Default("debug").
		Envar("LOG_LEVEL").
		EnumVar(&cfg.Log.Level, "debug", "info", "warning", "error", "fatal", "panic")
	kingpin.Flag("log-format", "Log format").
		Default("text").
		Envar("LOG_FORMAT").
		EnumVar(&cfg.Log.Format, "text", "json")

	kingpin.Flag("server-port", "HTTP server port").
		Default("5000").
		Envar("SERVER_PORT").
		StringVar(&cfg.Server.Port)

	kingpin.Flag("postgres-dsn", "DSN for postgres connection").
		Default("postgres://postgres:password@127.0.0.1:5643/postgres?sslmode=disable").
		Envar("POSTGRES_DSN").
		StringVar(&cfg.Postgres.DSN)

	kingpin.Flag("cors-allowed-origins", "AllowedOrigins for cors").
		Default("https://example.com").
		Envar("CORS_ALLOWED_ORIGINS").
		StringVar(&cfg.Cors.AllowedOrigins)

	kingpin.Flag("session-cookies-domain", "Domain for session cookies").
		Default(".example.com").
		Envar("SESSION_COOKIES_DOMAIN").
		StringVar(&cfg.Auth.SessionCookiesDomain)

	kingpin.Flag("access-token-secret-key", "Secret key for access token encoding/decoding").
		Default("d9049acad9744ced9c421d07c705b80e").
		Envar("ACCESS_TOKEN_SECRET_KEY").
		StringVar(&cfg.Auth.AccessTokenSecretKey)

	kingpin.Flag("refresh-token-secret-key", "Secret key for refresh token encoding/decoding").
		Default("e08b508c70d124c9dec74497daca9409").
		Envar("REFRESH_TOKEN_SECRET_KEY").
		StringVar(&cfg.Auth.RefreshTokenSecretKey)

	kingpin.Flag("google-auth-client-id", "GoogleAuthClientID for google OAuth").
		Default("long-hash.apps.googleusercontent.com").
		Envar("GOOGLE_AUTH_CLIENT_ID").
		StringVar(&cfg.Auth.GoogleAuthClientID)

	kingpin.Flag("google-auth-client-secret", "GoogleAuthClientSecret for google OAuth").
		Default("long_hash").
		Envar("GOOGLE_AUTH_CLIENT_SECRET").
		StringVar(&cfg.Auth.GoogleAuthClientSecret)

	kingpin.Flag("google-auth-callback-redirect-url", "URL for Google social auth callback redirect").
		Default("https://api.example.com/api/v1.0/auth/google/callback").
		Envar("GOOGLE_AUTH_CALLBACK_REDIRECT_URL").
		StringVar(&cfg.Auth.GoogleAuthCallbackRedirectURL)

	kingpin.Flag("facebook-auth-client-id", "FacebookAuthClientID for facebook OAuth").
		Default("hash").
		Envar("FACEBOOK_AUTH_CLIENT_ID").
		StringVar(&cfg.Auth.FacebookAuthClientID)

	kingpin.Flag("facebook-auth-client-secret", "FacebookAuthClientSecret for facebook OAuth").
		Default("hash").
		Envar("FACEBOOK_AUTH_CLIENT_SECRET").
		StringVar(&cfg.Auth.FacebookAuthClientSecret)

	kingpin.Flag("facebook-auth-callback-redirect-url", "URL for Facebook social auth callback redirect").
		Default("https://api.example.com/api/v1.0/auth/facebook/callback").
		Envar("FACEBOOK_AUTH_CALLBACK_REDIRECT_URL").
		StringVar(&cfg.Auth.FacebookAuthCallbackRedirectURL)

	kingpin.Flag("frontend-app-url", "URL of the frontend app").
		Default("https://example.com").
		Envar("FRONTEND_APP_URL").
		StringVar(&cfg.FrontendAppURL)

	kingpin.Flag("need-to-check-client", "Need to check trusted client header or search bot cookie").
		Default("true").
		Envar("NEED_TO_CHECK_CLIENT").
		BoolVar(&cfg.NeedToCheckClient)

	kingpin.Parse()
	return cfg
}
