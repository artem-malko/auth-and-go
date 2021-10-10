package manager

import (
	"testing"

	"github.com/artem-malko/auth-and-go/managers/user"

	. "github.com/artem-malko/auth-and-go/forks/goblin"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/services/account"
)

func TestDeleteUserByID(t *testing.T) {
	g := Goblin(t)
	env := &testEnv{}

	g.Describe("DeleteUserByID", func() {
		g.BeforeEach(func() {
			beforeEach(env)
		})

		g.AfterEach(func() {
			env.testingUserManager.db.Close()
		})

		g.It("Should delete User by ID without errors", func() {
			env.dbMock.ExpectBegin()
			env.dbMock.ExpectCommit()

			testUUID := uuid.MustParse("baa280a2-1421-42bf-8a99-26d7f3b990e9")

			env.mockedSessionService.On("DeleteAllSessionsByAccountID", testUUID, mock.Anything).Return(nil)
			env.mockedIdentityService.On("DeleteIdentitiesByAccountID", testUUID, mock.Anything).Return(nil)
			env.mockedAccountService.On("DeactivateAccountByID", testUUID, mock.Anything).Return(nil)

			res := env.testingUserManager.DeleteUserByID(testUUID)

			g.Assert(res).Equal(nil)

			env.mockedAccountService.AssertExpectations(t)
			env.mockedIdentityService.AssertExpectations(t)
			env.mockedSessionService.AssertExpectations(t)
		})

		g.It("Should return error, because of DeactivateAccountByID fail", func() {
			env.dbMock.ExpectBegin()
			env.dbMock.ExpectRollback()

			testUUID := uuid.MustParse("baa280a2-1421-42bf-8a99-26d7f3b990e9")

			env.mockedSessionService.On("DeleteAllSessionsByAccountID", testUUID, mock.Anything).Return(nil)
			env.mockedIdentityService.On("DeleteIdentitiesByAccountID", testUUID, mock.Anything).Return(nil)
			env.mockedAccountService.On("DeactivateAccountByID", testUUID, mock.Anything).Return(account.ErrNoAccountsUpdated)

			res := env.testingUserManager.DeleteUserByID(testUUID)

			g.Assert(res).Equal(user.ErrUserIsNotUpdated)

			env.mockedAccountService.AssertExpectations(t)
			env.mockedIdentityService.AssertExpectations(t)
			env.mockedSessionService.AssertExpectations(t)
		})
	})
}
