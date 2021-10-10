package service

import (
	"github.com/artem-malko/auth-and-go/services/token"
)

type tokenService struct {
	tokenRepository token.Repository
}

func New(tokenRepository token.Repository) token.Service {
	return &tokenService{
		tokenRepository: tokenRepository,
	}
}
