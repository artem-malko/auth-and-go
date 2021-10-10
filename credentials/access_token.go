package credentials

import (
	"encoding/base64"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/infrastructure/crypto"
	"github.com/artem-malko/auth-and-go/infrastructure/random"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/pkg/errors"
)

type AccessTokenInfo struct {
	AccessToken uuid.UUID
	AccountID   uuid.UUID
	SessionID   uuid.UUID
}

func CreateAccessToken(session models.Session, tokenSecretKey []byte) (string, error) {
	accessTokenPayload := random.RandStringRunes(5)

	accessTokenWithoutExpiresDate, err := crypto.Encrypt(
		strings.Join([]string{
			accessTokenPayload,
			session.AccessToken.String(),
			session.AccountID.String(),
			session.ID.String(),
		}, tokenValuesSeparator),
		tokenSecretKey,
	)

	if err != nil {
		return "", errors.Wrap(err, "access token generation error")
	}

	accessToken := base64.URLEncoding.EncodeToString([]byte(strings.Join([]string{
		accessTokenWithoutExpiresDate,
		session.AccessTokenExpiresDate.UTC().Format(time.RFC3339),
	}, tokenValuesSeparator)))

	return accessToken, nil
}

func ParseRawAccessToken(rawAccessToken string, decodeKey []byte) (*AccessTokenInfo, error) {
	unescapedAccessTokenString, err := url.PathUnescape(rawAccessToken)

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawAccessToken string error, unescape error")
	}

	decodedAccessTokenString, err := base64.URLEncoding.DecodeString(unescapedAccessTokenString)

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawAccessToken string error, decode error")
	}

	decodedAccessTokenStringParts := strings.Split(string(decodedAccessTokenString), tokenValuesSeparator)

	if len(decodedAccessTokenStringParts) != 2 {
		return nil, errors.Wrap(err, "ParseRawAccessToken string error, incorrect decoded parts count")
	}

	accessTokenExpiresDate, err := time.Parse(time.RFC3339, decodedAccessTokenStringParts[1])

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawAccessToken string error, parse expires date error")
	}

	duration := time.Since(accessTokenExpiresDate)

	if duration.Seconds() > 0 {
		return nil, errors.Wrap(err, "ParseRawAccessToken string error, access token expired")
	}

	decryptedAccessTokenString, err := crypto.Decrypt(decodedAccessTokenStringParts[0], decodeKey)

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawAccessToken string error, decrypt error")
	}

	decryptedAccessTokenStringParts := strings.Split(decryptedAccessTokenString, tokenValuesSeparator)

	if len(decryptedAccessTokenStringParts) != 4 {
		return nil, errors.New("ParseRawAccessToken string error, incorrect parts count")
	}

	accessToken, err := uuid.Parse(decryptedAccessTokenStringParts[1])

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawAccessToken string error, accessToken parsing error")
	}

	accountID, err := uuid.Parse(decryptedAccessTokenStringParts[2])

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawAccessToken string error, accountID parsing error")
	}

	sessionID, err := uuid.Parse(decryptedAccessTokenStringParts[3])

	if err != nil {
		return nil, errors.Wrap(err, "ParseRawAccessToken string error, sessionID parsing error")
	}

	return &AccessTokenInfo{
		AccessToken: accessToken,
		AccountID:   accountID,
		SessionID:   sessionID,
	}, nil
}
