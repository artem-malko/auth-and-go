package request

import (
	"mime/multipart"
	"net/http"

	"github.com/golang/gddo/httputil/header"

	"github.com/apex/log"
)

type MultiPartBinder interface {
	Bind(r *http.Request, form *multipart.Form) error
}

func ParseMultipartFormDataBody(
	r *http.Request,
	binder MultiPartBinder,
	caller string,
	logger func(r *http.Request) log.Interface,
) *ParseBodyError {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "multipart/form-data" {
			return &ParseBodyError{
				Status:  http.StatusUnsupportedMediaType,
				Message: "Content-Type header is not multipart/form-data",
			}
		}
	}

	parseErr := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if parseErr != nil {
		// @TODO Add more info about err, add source?
		logger(r).
			WithField("method", caller).
			WithField("code", http.StatusInternalServerError).
			Error(parseErr.Error())

		return &ParseBodyError{
			Status:  http.StatusBadRequest,
			Message: "Failed to parse multipart/form-data",
		}
	}

	if r.MultipartForm == nil {
		return &ParseBodyError{
			Status:  http.StatusBadRequest,
			Message: "Request body must not be empty",
		}
	}

	err := binder.Bind(r, r.MultipartForm)

	if err != nil {
		return &ParseBodyError{Status: http.StatusBadRequest, Message: err.Error()}
	}

	return nil
}
