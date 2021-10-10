package managers

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/services"
	"github.com/pkg/errors"

	"github.com/apex/log"

	userManager "github.com/artem-malko/auth-and-go/managers/user/manager"

	"github.com/artem-malko/auth-and-go/managers/user"
)

type Managers struct {
	UserManager user.Manager
}

type NewManagersParams struct {
	DB                    *sql.DB
	Logger                log.Interface
	AccessTokenSecretKey  string
	RefreshTokenSecretKey string
}

func New(params NewManagersParams) (*Managers, error) {
	servicesInstance, err := services.New(params.DB, params.Logger)

	if err != nil {
		return nil, errors.Wrap(err, "Services creation failed")
	}

	userManagerInstance := userManager.New(
		params.DB,
		params.Logger,
		servicesInstance.AccountService,
		servicesInstance.IdentityService,
		servicesInstance.SessionService,
		servicesInstance.TokenService,
		servicesInstance.MailerService,
		params.AccessTokenSecretKey,
		params.RefreshTokenSecretKey,
	)

	return &Managers{
		UserManager: userManagerInstance,
	}, nil
}
