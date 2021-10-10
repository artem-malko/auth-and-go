package logger

import (
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/pkg/errors"
)

func New(format, level string) (log.Interface, error) {
	var handler log.Handler

	switch format {
	case "text":
		handler = text.New(os.Stdout)
	case "json":
		handler = json.New(os.Stdout)
	default:
		return nil, errors.Errorf("logger: invalid format (%s)", format)
	}

	logLevel, logLevelParseErr := log.ParseLevel(level)

	if logLevelParseErr != nil {
		return nil, errors.Wrapf(logLevelParseErr, "logger: invalid level (%s)", level)
	}

	return &log.Logger{
		Handler: handler,
		Level:   logLevel,
	}, nil
}

func NewForTests() log.Interface {
	return &log.Logger{
		Handler: text.New(ioutil.Discard),
		Level:   log.DebugLevel,
	}
}
