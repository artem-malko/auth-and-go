package manager

import (
	"testing"

	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/models"

	"github.com/stretchr/testify/mock"

	"github.com/artem-malko/auth-and-go/services/identity"

	. "github.com/artem-malko/auth-and-go/forks/goblin"
)

func TestCreateUserWithEmail(t *testing.T) {
	g := Goblin(t)
	env := &testEnv{}

	g.Describe("CreateUserWithEmail", func() {
		g.BeforeEach(func() {
			beforeEach(env)
		})

		g.AfterEach(func() {
			env.testingUserManager.db.Close()
		})

		g.It("Should create new email identity and account for absolutely new user without errors", func() {
			env.dbMock.ExpectBegin()
			env.dbMock.ExpectCommit()

			accountID := uuid.New()
			identityID := uuid.New()
			tokenID := uuid.New()
			email := "test@test.com"
			password := "password"
			env.mockedIdentityService.
				On("GetIdentitiesByEmail", email).
				Return(nil, identity.ErrNoIdentityFound)
			env.mockedAccountService.
				On("CreateAccount", mock.Anything, email, models.AccountStatusUnconfirmed).
				Return(&models.Account{
					ID:            accountID,
					AccountStatus: models.AccountStatusUnconfirmed,
				}, nil)
			env.mockedIdentityService.
				On("CreateEmailIdentity", accountID, email, password, models.IdentityStatusUnconfirmed, mock.Anything).
				Return(&models.Identity{
					ID:           identityID,
					AccountID:    accountID,
					Email:        email,
					IdentityType: models.IdentityTypeEmail,
				}, nil)
			env.mockedTokenService.
				On("Create", models.TokenTypeRegistrationConfirmation, models.ClientIDWEB, accountID, identityID, mock.Anything).
				Return(&models.Token{ID: tokenID}, nil)
			env.mockedMailerService.
				On("SendRegistrationConfirmationEmail", email, tokenID).
				Return(nil)

			res := env.testingUserManager.CreateUserWithEmail(email, password, models.ClientIDWEB)

			g.Assert(res).Equal(nil)

			env.mockedAccountService.AssertExpectations(g)
			env.mockedIdentityService.AssertExpectations(g)
			env.mockedSessionService.AssertExpectations(g)
			env.mockedTokenService.AssertExpectations(g)
			env.mockedMailerService.AssertExpectations(g)
		})
	})
}
