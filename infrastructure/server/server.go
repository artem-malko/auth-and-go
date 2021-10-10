package server

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// События генерируемые до/после совершения действий над веб сервером
const (
	EventBeforeStop  = "web:before.stop"
	EventAfterStop   = "web:after.stop"
	EventBeforeStart = "web:before.start"
)

// Server с публикацией событйи о своем состоянии в Event Bus
type Server struct {
	server *http.Server
	pub    eventPublisher
}

type eventPublisher interface {
	Publish(topic string, args ...interface{})
}

// New инициализация севрера с настройками и зависимостями
func New(addr string, h http.Handler, p eventPublisher) *Server {
	return &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      h,
			ReadTimeout:  10000 * time.Second,
			WriteTimeout: 10000 * time.Second,
		},
		pub: p,
	}
}

// Start web сервера
func (s *Server) Start() error {
	errChan := make(chan error, 1)
	go func() {
		s.pub.Publish(EventBeforeStart)

		errChan <- s.server.ListenAndServe()
	}()
	if err := <-errChan; err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrap(err, "web error")
		}
	}
	return nil
}

// Stop web сервера c плавным завершением всех входящих запросов
func (s *Server) Stop() error {
	s.pub.Publish(EventBeforeStop)

	if err := s.server.Shutdown(context.Background()); err != nil {
		return errors.Wrap(err, "could not shutdown web")
	}

	s.pub.Publish(EventAfterStop)
	return nil
}
