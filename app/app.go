package app

import "github.com/apex/log"

type App struct {
	log      log.Interface
	stopChan chan struct{}
}

func New() *App {
	return &App{}
}

func (a *App) Run() error {
	return nil
}
