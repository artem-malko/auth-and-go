package services

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/services/mailer"

	"github.com/artem-malko/auth-and-go/repositories"

	"github.com/artem-malko/auth-and-go/services/token"

	"github.com/artem-malko/auth-and-go/services/identity"

	"github.com/apex/log"
	"github.com/artem-malko/auth-and-go/services/account"

	"github.com/artem-malko/auth-and-go/services/session"

	accountService "github.com/artem-malko/auth-and-go/services/account/service"
	identityService "github.com/artem-malko/auth-and-go/services/identity/service"
	mailerService "github.com/artem-malko/auth-and-go/services/mailer/service"
	sessionService "github.com/artem-malko/auth-and-go/services/session/service"
	tokenService "github.com/artem-malko/auth-and-go/services/token/service"

	"github.com/pkg/errors"
)

type Services struct {
	SessionService  session.Service
	AccountService  account.Service
	IdentityService identity.Service
	TokenService    token.Service
	MailerService   mailer.Service
}

func New(db *sql.DB, logger log.Interface) (*Services, error) {
	repositoryInstances, err := repositories.New(db, logger)

	if err != nil {
		return nil, errors.Wrap(err, "services: repositories init error")
	}

	return &Services{
		SessionService:  sessionService.New(repositoryInstances.SessionRepository),
		AccountService:  accountService.New(repositoryInstances.AccountRepository),
		IdentityService: identityService.New(repositoryInstances.IdentityRepository),
		TokenService:    tokenService.New(repositoryInstances.TokenRepository),
		MailerService:   mailerService.New(),
	}, nil
}
