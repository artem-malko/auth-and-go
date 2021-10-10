package manager

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/apex/log"

	"github.com/artem-malko/auth-and-go/services/mailer"
	"github.com/artem-malko/auth-and-go/services/token"

	"github.com/artem-malko/auth-and-go/managers/user"

	"github.com/artem-malko/auth-and-go/services/account"
	"github.com/artem-malko/auth-and-go/services/identity"
	"github.com/artem-malko/auth-and-go/services/session"
)

type userManager struct {
	identityService       identity.Service
	accountService        account.Service
	sessionService        session.Service
	tokenService          token.Service
	mailerService         mailer.Service
	accessTokenSecretKey  []byte
	refreshTokenSecretKey []byte
	db                    *sql.DB
	logger                log.Interface
}

func New(
	db *sql.DB,
	logger log.Interface,
	accountService account.Service,
	identityService identity.Service,
	sessionService session.Service,
	tokenService token.Service,
	mailerService mailer.Service,
	accessTokenSecretKey string,
	refreshTokenSecretKey string,
) user.Manager {
	return &userManager{
		identityService:       identityService,
		accountService:        accountService,
		sessionService:        sessionService,
		tokenService:          tokenService,
		mailerService:         mailerService,
		db:                    db,
		accessTokenSecretKey:  []byte(accessTokenSecretKey),
		refreshTokenSecretKey: []byte(refreshTokenSecretKey),
		logger:                logger,
	}
}

func (m *userManager) runWithTransaction(method string, function func(tx *sql.Tx) error) error {
	return database.RunWithTransaction(m.db, method, function)
}
