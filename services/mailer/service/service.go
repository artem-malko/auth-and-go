package service

import (
	"bytes"
	"html/template"
	"net/smtp"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/services/mailer"
)

type mailerService struct {
	from string
	host string
	port string
	auth smtp.Auth
}

func New() mailer.Service {
	//
	// user we are authorizing as
	from := "artem.malko@gmail.com"

	// server we are authorized to send email through
	host := "smtp.gmail.com"

	return &mailerService{
		from: from,
		host: host,
		port: "587",
		auth: smtp.PlainAuth("", from, "qwe", host),
	}

}

func (s *mailerService) parseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)

	if err != nil {
		return "", errors.Wrap(err, "ParseTemplate ParseFiles error")
	}
	buf := new(bytes.Buffer)

	if err = t.Execute(buf, data); err != nil {
		return "", errors.Wrap(err, "ParseTemplate execute error")
	}

	return buf.String(), nil
}
