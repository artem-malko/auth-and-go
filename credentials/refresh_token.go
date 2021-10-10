package credentials

import (
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/infrastructure/crypto"
	"github.com/artem-malko/auth-and-go/infrastructure/random"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/pkg/errors"
)

type RefreshTokenInfo struct {
	RefreshToken uuid.UUID
}

func CreateRefreshToken(session models.Session, tokenSecretKey []byte) (string, error) {
	refreshTokenPayload := random.RandStringRunes(5)

	refreshToken, err := crypto.Encrypt(
		strings.Join([]string{
			refreshTokenPayload,
			session.RefreshTokenExpiresDate.UTC().Format(time.RFC3339),
			session.RefreshToken.String(),
		}, tokenValuesSeparator),
		tokenSecretKey,
	)

	if err != nil {
		return "", errors.Wrap(err, "refresh token generation error")
	}

	return refreshToken, nil
}

func ParseRawRefreshToken(rawRefreshToken string, decodeKey []byte) (*RefreshTokenInfo, error) {
	decodedRefreshTokenString, err := url.PathUnescape(rawRefreshToken)

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawRefreshToken string error, decode error")
	}

	decryptedRefreshTokenString, err := crypto.Decrypt(decodedRefreshTokenString, decodeKey)

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawRefreshToken string error, decrypt error")
	}

	decryptedRefreshTokenStringParts := strings.Split(decryptedRefreshTokenString, tokenValuesSeparator)

	if len(decryptedRefreshTokenStringParts) != 3 {
		return nil, errors.Wrap(err, "ParseRawRefreshToken string error, incorrect parts count")
	}

	refreshTokenExpiresDate, err := time.Parse(time.RFC3339, decryptedRefreshTokenStringParts[1])

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawRefreshToken string error, parse expires date error")
	}

	duration := time.Since(refreshTokenExpiresDate)

	if duration.Seconds() > 0 {
		return nil, errors.Wrap(err, "ParseRawRefreshToken string error, refresh token expired")
	}

	refreshToken, err := uuid.Parse(decryptedRefreshTokenStringParts[2])

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawRefreshToken string error, refreshToken parsing error")
	}

	return &RefreshTokenInfo{
		RefreshToken: refreshToken,
	}, nil
}
