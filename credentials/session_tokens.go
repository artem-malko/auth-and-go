package credentials

import (
	"github.com/artem-malko/auth-and-go/models"
)

const tokenValuesSeparator = ","

func CreateSessionTokens(
	session models.Session,
	accessTokenSecretKey, refreshTokenSecretKey []byte,
) (*models.SessionTokens, error) {
	accessToken, err := CreateAccessToken(session, accessTokenSecretKey)

	if err != nil {
		return nil, err
	}

	refreshToken, err := CreateRefreshToken(session, refreshTokenSecretKey)

	if err != nil {
		return nil, err
	}

	return &models.SessionTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
