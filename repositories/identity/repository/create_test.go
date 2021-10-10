package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/repositories/identity"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/artem-malko/auth-and-go/forks/goblin"
	"github.com/artem-malko/auth-and-go/models"
)

func TestCreateEmailIdentity(t *testing.T) {
	g := Goblin(t)
	var dbMock sqlmock.Sqlmock
	var db *sql.DB
	var testingRepository identity.Repository

	g.Describe("CreateEmailIdentity", func() {
		g.BeforeEach(func() {
			db, dbMock, _ = sqlmock.New()
			dbMock.ExpectExec("ANALYZE identities").WillReturnResult(sqlmock.NewResult(1, 1))
			testingRepository, _ = New()
		})

		g.AfterEach(func() {
			db.Close()
		})

		g.It("Create new email identity in transaction without errors", func() {
			dbMock.ExpectBegin()
			executor, _ := db.Begin()

			testingIdentity := models.Identity{
				ID:             uuid.New(),
				AccountID:      uuid.New(),
				IdentityType:   models.IdentityTypeEmail,
				IdentityStatus: models.IdentityStatusUnconfirmed,
				Email:          "example@test.com",
				PasswordHash:   "hash",
			}
			dbMock.ExpectQuery("INSERT INTO identities").
				WithArgs(
					testingIdentity.ID,
					testingIdentity.AccountID,
					testingIdentity.IdentityType,
					testingIdentity.IdentityStatus,
					testingIdentity.Email,
					testingIdentity.PasswordHash,
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(testingIdentity.ID, time.Time{}, time.Time{}),
				)

			createdIdentity, err := testingRepository.CreateEmailIdentity(executor, testingIdentity)

			assert.Equal(g, nil, err)
			assert.Equal(g, testingIdentity.ID, createdIdentity.ID)
		})

		g.It("Create new email identity in transaction with not uniq email error", func() {
			dbMock.ExpectBegin()
			executor, _ := db.Begin()

			testingIdentity := models.Identity{
				Email:        "example@test.com",
				PasswordHash: "hash",
			}

			dbMock.ExpectQuery("INSERT INTO identities").
				WithArgs(
					testingIdentity.ID,
					testingIdentity.AccountID,
					testingIdentity.IdentityType,
					testingIdentity.IdentityStatus,
					testingIdentity.Email,
					testingIdentity.PasswordHash,
				).
				WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation, ConstraintName: emailUniqIndexName})

			_, err := testingRepository.CreateEmailIdentity(executor, testingIdentity)

			assert.Equal(g, identity.ErrRepositoryEmailConstraint, err)
		})

		g.It("Create new email identity in transaction with any pgxErr error, but not emailUniq constraint", func() {
			dbMock.ExpectBegin()
			executor, _ := db.Begin()

			testingIdentity := models.Identity{
				Email:        "example@test.com",
				PasswordHash: "hash",
			}

			dbMock.ExpectQuery("INSERT INTO identities").
				WithArgs(
					testingIdentity.ID,
					testingIdentity.AccountID,
					testingIdentity.IdentityType,
					testingIdentity.IdentityStatus,
					testingIdentity.Email,
					testingIdentity.PasswordHash,
				).
				WillReturnError(&pgconn.PgError{Code: "23030"})

			_, err := testingRepository.CreateEmailIdentity(executor, testingIdentity)

			assert.Equal(g, &pgconn.PgError{Code: "23030"}, errors.Cause(err))
		})
	})
}

func TestCreateSocialIdentity(t *testing.T) {
	g := Goblin(t)
	var dbMock sqlmock.Sqlmock
	var db *sql.DB
	var testingRepository identity.Repository

	g.Describe("CreateSocialIdentity", func() {
		g.BeforeEach(func() {
			db, dbMock, _ = sqlmock.New()
			dbMock.ExpectExec("ANALYZE identities").WillReturnResult(sqlmock.NewResult(1, 1))
			testingRepository, _ = New()
		})

		g.AfterEach(func() {
			db.Close()
		})

		g.It("Create new google identity in transaction without errors", func() {
			dbMock.ExpectBegin()
			executor, _ := db.Begin()

			testingIdentity := models.Identity{
				ID:             uuid.New(),
				AccountID:      uuid.New(),
				IdentityType:   models.IdentityTypeGoogle,
				IdentityStatus: models.IdentityStatusConfirmed,
				Email:          "example@test.com",
				GoogleSocialID: "GoogleSocialID",
			}
			dbMock.ExpectQuery("INSERT INTO identities").
				WithArgs(
					testingIdentity.ID,
					testingIdentity.AccountID,
					testingIdentity.IdentityType,
					testingIdentity.IdentityStatus,
					testingIdentity.Email,
					testingIdentity.GoogleSocialID,
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(testingIdentity.ID, time.Time{}, time.Time{}),
				)

			createdIdentity, err := testingRepository.CreateSocialIdentity(executor, testingIdentity)

			assert.Equal(g, nil, err)
			assert.Equal(g, createdIdentity.ID, testingIdentity.ID)
			assert.Equal(g, createdIdentity.GoogleSocialID, testingIdentity.GoogleSocialID)
		})

		g.It("Create new facebook identity in transaction without errors", func() {
			dbMock.ExpectBegin()
			executor, _ := db.Begin()

			testingIdentity := models.Identity{
				ID:               uuid.New(),
				AccountID:        uuid.New(),
				IdentityType:     models.IdentityTypeFacebook,
				IdentityStatus:   models.IdentityStatusConfirmed,
				Email:            "example@test.com",
				FacebookSocialID: "FacebookSocialID",
			}
			dbMock.ExpectQuery("INSERT INTO identities").
				WithArgs(
					testingIdentity.ID,
					testingIdentity.AccountID,
					testingIdentity.IdentityType,
					testingIdentity.IdentityStatus,
					testingIdentity.Email,
					testingIdentity.FacebookSocialID,
				).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(testingIdentity.ID, time.Time{}, time.Time{}),
				)

			createdIdentity, err := testingRepository.CreateSocialIdentity(executor, testingIdentity)

			assert.Equal(g, nil, err)
			assert.Equal(g, createdIdentity.ID, testingIdentity.ID)
			assert.Equal(g, createdIdentity.FacebookSocialID, testingIdentity.FacebookSocialID)
		})

		g.It("Creating new unknown social identity returns ErrRepositoryUnknownSocialNetworkType error", func() {
			testingIdentity := models.Identity{
				ID:             uuid.New(),
				AccountID:      uuid.New(),
				IdentityType:   "unknown",
				IdentityStatus: models.IdentityStatusConfirmed,
				Email:          "example@test.com",
			}

			_, err := testingRepository.CreateSocialIdentity(nil, testingIdentity)

			assert.Equal(g, identity.ErrRepositoryUnknownSocialNetworkType, err)
		})

		g.It("Creating social identity with existed GoogleSocialID returns pgx constraint error", func() {
			testingIdentity := models.Identity{
				ID:             uuid.New(),
				AccountID:      uuid.New(),
				IdentityType:   models.IdentityTypeGoogle,
				IdentityStatus: models.IdentityStatusConfirmed,
				Email:          "example@test.com",
				GoogleSocialID: "GoogleSocialID",
			}
			dbMock.ExpectQuery("INSERT INTO identities").
				WithArgs(
					testingIdentity.ID,
					testingIdentity.AccountID,
					testingIdentity.IdentityType,
					testingIdentity.IdentityStatus,
					testingIdentity.Email,
					testingIdentity.GoogleSocialID,
				).
				WillReturnError(&pgconn.PgError{
					Code:           pgerrcode.UniqueViolation,
					ConstraintName: googleSocialIDUniqIndexName,
				})

			_, err := testingRepository.CreateSocialIdentity(nil, testingIdentity)

			assert.Equal(g, identity.ErrRepositorySocialGoogleConstraint, err)
		})

		g.It("Creating social identity with existed FacebookSocialID returns pgx constraint error", func() {
			testingIdentity := models.Identity{
				ID:               uuid.New(),
				AccountID:        uuid.New(),
				IdentityType:     models.IdentityTypeFacebook,
				IdentityStatus:   models.IdentityStatusConfirmed,
				Email:            "example@test.com",
				FacebookSocialID: "FacebookSocialID",
			}
			dbMock.ExpectQuery("INSERT INTO identities").
				WithArgs(
					testingIdentity.ID,
					testingIdentity.AccountID,
					testingIdentity.IdentityType,
					testingIdentity.IdentityStatus,
					testingIdentity.Email,
					testingIdentity.FacebookSocialID,
				).
				WillReturnError(
					&pgconn.PgError{
						Code:           pgerrcode.UniqueViolation,
						ConstraintName: facebookSocialIDUniqIndexName,
					})

			_, err := testingRepository.CreateSocialIdentity(nil, testingIdentity)

			assert.Equal(g, identity.ErrRepositorySocialFacebookConstraint, err)
		})
	})
}
