package service

import (
	"github.com/artem-malko/auth-and-go/services/session"
)

type sessionService struct {
	sessionRepository session.Repository
}

func New(sessionRepository session.Repository) session.Service {
	return &sessionService{
		sessionRepository: sessionRepository,
	}
}
