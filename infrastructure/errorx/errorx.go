package errorx

import (
	"github.com/artem-malko/auth-and-go/infrastructure/caller"
	"github.com/pkg/errors"
)

type ErrorHandler func(err error, message ...string) error

func CreateErrorHandler(prefix string) ErrorHandler {
	return func(err error, message ...string) error {
		callerName := "unknown " + prefix + " caller"

		if name, ok := caller.GetCaller(1); ok == true {
			callerName = name
		}

		defMessage := prefix + " Caller: " + callerName

		if len(message) != 0 {
			return errors.Wrap(err, defMessage+" msg: "+message[0])
		}

		return errors.Wrap(err, defMessage)
	}
}
