package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/artem-malko/auth-and-go/forks/goblin"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func TestGetUserByID(t *testing.T) {
	g := Goblin(t)
	env := &testEnv{}

	g.Describe("GetUserByID", func() {
		g.BeforeEach(func() {
			beforeEach(env)
		})

		g.AfterEach(func() {
			env.testingUserManager.db.Close()
		})

		g.It("Should return user without identities and errors", func() {
			userID := uuid.MustParse("baa280a2-1421-42bf-8a99-26d7f3b990e9")
			// setup expectations
			env.mockedAccountService.On("GetAccountByID", userID).Return(&models.Account{
				ID: userID,
			}, nil)

			userByID, _ := env.testingUserManager.GetUserByID(userID)
			var userIdentities []models.UserIdentity

			// assert equality
			assert.Equal(g, userByID.ID, userID, "they should be equal")
			assert.Equal(g, userIdentities, userByID.Identities, "Identities should be nil slice")

			// assert that the expectations were met
			env.mockedAccountService.AssertExpectations(g)
		})
	})
}
