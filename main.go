package main

import (
	"log"

	"github.com/artem-malko/auth-and-go/app"
	"github.com/pkg/errors"
)

func main() {
	application := app.New()
	err := application.Run()

	if err != nil {
		log.Fatal(errors.Wrapf(err, "App init error"))
	}
}
