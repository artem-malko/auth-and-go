package repository

import (
	"database/sql"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/caller"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/artem-malko/auth-and-go/infrastructure/convert"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/repositories/account"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

const accountsTableName = "accounts"
const repositoryPrefix = accountsTableName + " repository: "
const idFieldName = "id"
const accountTypeFieldName = "account_type"
const accountStatusFieldName = "account_status"
const accountNameFieldName = "account_name"
const profileFieldName = "profile"
const settingsFieldName = "settings"
const lastIPFieldName = "last_ip"
const lastLoginFieldName = "last_login"
const createdAtFieldName = "created_at"
const updatedAtFieldName = "updated_at"

const accountNameUniqIndexName = "accounts_account_name_active_account_uniq_idx"

type dbProfile struct {
	FirstName   null.String         `json:"first_name"`
	LastName    null.String         `json:"last_name"`
	Country     null.String         `json:"country"`
	City        null.String         `json:"city"`
	Gender      null.String         `json:"gender"`
	Description null.String         `json:"description"`
	Birthday    null.Time           `json:"birthday"`
	SocialLinks []models.SocialLink `json:"social_links"`
	Interests   []string            `json:"interests"`
}

type dbEmailNotification struct {
	Email string                              `json:"email"`
	Scope []models.EmailNotificationScopeType `json:"scope"`
}

type dbNotifications struct {
	Email dbEmailNotification `json:"email"`
}

type dbSettings struct {
	Notifications dbNotifications `json:"notifications"`
}

// User for DB
type dbAccount struct {
	ID            uuid.UUID
	AccountStatus models.AccountStatus
	AccountType   models.AccountType
	AccountName   string
	Profile       dbProfile
	LastIP        string
	LastLogin     time.Time
	Settings      dbSettings
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (a *dbAccount) GetBaseAccount() models.Account {
	return models.Account{
		ID:            a.ID,
		AccountStatus: a.AccountStatus,
		AccountType:   a.AccountType,
		AccountName:   a.AccountName,
		Profile: models.Profile{
			FirstName:   a.Profile.FirstName.String,
			LastName:    a.Profile.LastName.String,
			Country:     a.Profile.Country.String,
			City:        a.Profile.City.String,
			Gender:      a.Profile.Gender.String,
			Description: a.Profile.Description.String,
			Birthday:    a.Profile.Birthday.Time,
			SocialLinks: a.Profile.SocialLinks,
			Interests:   a.Profile.Interests,
		},
		Settings: models.Settings{
			Notifications: models.Notifications{
				Email: models.EmailNotifications{
					Email: a.Settings.Notifications.Email.Email,
					Scope: a.Settings.Notifications.Email.Scope,
				},
			},
		},
		LastIP:    a.LastIP,
		LastLogin: a.LastLogin,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func createDBAccount(a models.Account) dbAccount {
	return dbAccount{
		ID:            a.ID,
		AccountStatus: a.AccountStatus,
		AccountType:   a.AccountType,
		AccountName:   a.AccountName,
		Profile: dbProfile{
			FirstName:   convert.NewSQLNullString(a.Profile.FirstName),
			LastName:    convert.NewSQLNullString(a.Profile.LastName),
			Country:     convert.NewSQLNullString(a.Profile.Country),
			City:        convert.NewSQLNullString(a.Profile.City),
			Gender:      convert.NewSQLNullString(a.Profile.Gender),
			Birthday:    convert.NewSQLNullTime(a.Profile.Birthday),
			Description: convert.NewSQLNullString(a.Profile.Description),
			SocialLinks: a.Profile.SocialLinks,
			Interests:   a.Profile.Interests,
		},
		LastIP:    a.LastIP,
		LastLogin: a.LastLogin,
		Settings: dbSettings{
			Notifications: dbNotifications{
				Email: dbEmailNotification{
					Email: a.Settings.Notifications.Email.Email,
					Scope: a.Settings.Notifications.Email.Scope,
				},
			},
		},
		UpdatedAt: a.UpdatedAt,
		CreatedAt: a.CreatedAt,
	}
}

// New creates postgresRepository
func New() (account.Repository, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &repository{
		psql: psql,
	}, nil
}

type repository struct {
	psql sq.StatementBuilderType
}

func getCaller() string {
	// Remove getCaller from stack
	if callerName, ok := caller.GetCaller(2); ok == true {
		return callerName
	}

	return repositoryPrefix + "unknown caller"
}

func (r *repository) exec(executor database.QueryExecutor, queryBuilder database.QueryBuilder) error {
	return database.ExecWithCaller(
		getCaller(),
		executor,
		queryBuilder,
	)
}

func (r *repository) queryRow(executor database.QueryExecutor, queryBuilder database.QueryBuilder) (*sql.Row, error) {
	return database.QueryRowWithCaller(
		getCaller(),
		executor,
		queryBuilder,
	)
}

func (r *repository) query(
	executor database.QueryExecutor,
	queryBuilder database.QueryBuilder,
	scanFunc func(rows *sql.Rows) error,
) error {
	return database.QueryWithCaller(
		getCaller(),
		executor,
		queryBuilder,
		scanFunc,
	)
}
