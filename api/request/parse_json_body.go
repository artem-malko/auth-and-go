package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"

	"github.com/apex/log"
)

type JSONBinder interface {
	Bind(r *http.Request) error
}

type ParseBodyError struct {
	Status  int
	Message string
}

func ParseJSONBody(r *http.Request, dst JSONBinder, caller string, logger func(r *http.Request) log.Interface) *ParseBodyError {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			return &ParseBodyError{
				Status:  http.StatusUnsupportedMediaType,
				Message: "Content-Type header is not application/json",
			}
		}
	}

	// @TODO maybe later
	//r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			message := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &ParseBodyError{
				Status:  http.StatusBadRequest,
				Message: message,
			}

		case errors.Is(err, io.ErrUnexpectedEOF):
			message := "Request body contains badly-formed JSON"
			return &ParseBodyError{Status: http.StatusBadRequest, Message: message}

		case errors.As(err, &unmarshalTypeError):
			message := fmt.Sprintf(
				"Request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field,
				unmarshalTypeError.Offset,
			)
			return &ParseBodyError{Status: http.StatusBadRequest, Message: message}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			message := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &ParseBodyError{Status: http.StatusBadRequest, Message: message}

		case errors.Is(err, io.EOF):
			message := "Request body must not be empty"
			return &ParseBodyError{Status: http.StatusBadRequest, Message: message}

		//case err.Error() == "http: request body too large":
		//	Message := "Request body must not be larger than 1MB"
		//	return &ParseBodyError{Status: http.StatusRequestEntityTooLarge, Message: Message}

		default:
			logger(r).
				WithField("method", caller).
				WithField("code", http.StatusInternalServerError).
				Error(err.Error())
			return &ParseBodyError{Status: http.StatusInternalServerError, Message: "Internal server error"}
		}
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		message := "Request body must only contain a single JSON object"
		return &ParseBodyError{Status: http.StatusBadRequest, Message: message}
	}

	err = dst.Bind(r)

	if err != nil {
		return &ParseBodyError{Status: http.StatusBadRequest, Message: err.Error()}
	}

	return nil
}
