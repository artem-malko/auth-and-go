package manager

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artem-malko/auth-and-go/services/account"
	"github.com/artem-malko/auth-and-go/services/identity"
	"github.com/artem-malko/auth-and-go/services/mailer"
	"github.com/artem-malko/auth-and-go/services/session"
	"github.com/artem-malko/auth-and-go/services/token"
)

type testEnv struct {
	testingUserManager    *userManager
	mockedAccountService  *account.MockService
	mockedIdentityService *identity.MockService
	mockedSessionService  *session.MockService
	mockedTokenService    *token.MockService
	mockedMailerService   *mailer.MockService
	dbMock                sqlmock.Sqlmock
	db                    *sql.DB
}

func beforeEach(env *testEnv) {
	db, dbMock, _ := sqlmock.New()
	env.db = db
	env.dbMock = dbMock
	env.mockedAccountService = new(account.MockService)
	env.mockedIdentityService = new(identity.MockService)
	env.mockedSessionService = new(session.MockService)
	env.mockedTokenService = new(token.MockService)
	env.mockedMailerService = new(mailer.MockService)
	env.testingUserManager = &userManager{
		accountService:  env.mockedAccountService,
		identityService: env.mockedIdentityService,
		sessionService:  env.mockedSessionService,
		tokenService:    env.mockedTokenService,
		mailerService:   env.mockedMailerService,
		db:              db,
	}
}
