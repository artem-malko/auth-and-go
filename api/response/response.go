package response

import (
	"encoding/json"
	"net/http"

	"github.com/artem-malko/auth-and-go/infrastructure/caller"

	"github.com/artem-malko/auth-and-go/api/context_utils"

	"github.com/apex/log"
)

const contentTypeJSON = "application/json; charset=UTF-8"

// ErrorResponse represents error response
type ErrorResponse struct {
	ErrorData ErrorData `json:"error"`
}

// ErrorData represents error in error response
type ErrorData struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// OK response with 200 status code and data in payload
func OK(w http.ResponseWriter, data interface{}) {
	jsonResponse(w, SuccessResponse{
		Data: data,
	})
}

// OKWithoutContent response with 204 status code
func OKWithoutContent(w http.ResponseWriter) {
	// We need to set Content-Type here, cause we can not do it after WriteHeader
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(204)
}

// OKWithResponseCode response with 2xx status code
func OKWithResponseCode(w http.ResponseWriter, statusCode int, data interface{}) {
	// We need to set Content-Type here, cause we can not do it after WriteHeader
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(statusCode)
	jsonResponse(w, SuccessResponse{
		Data: data,
	})
}

// InternalServerError replies to the request with an HTTP 500 Internal Server Error error
func InternalServerError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "Internal server error")
}

// CreateBoundLogger create error response generator
// Method will log error
func CreateBoundLogger(logger log.Interface) func(r *http.Request) log.Interface {
	return func(r *http.Request) log.Interface {
		ctx := r.Context()
		requestID := context_utils.GetRequestID(ctx)
		callerName := "Unknown"

		if name, ok := caller.GetCaller(2); ok == true {
			callerName = name
		}

		return logger.
			WithField("caller", callerName).
			WithField("request_id", requestID).
			WithField("source", "router")
	}
}

// NotFound replies to the request with an HTTP 404 Not Found error
func NotFound(w http.ResponseWriter) {
	Error(w, http.StatusNotFound, "Not Found")
}

// BadRequest replies to the request with an HTTP 404 Given Param is not valid
func BadRequest(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, "Given Params are not valid")
}

// NotAllowed replies to the request with an HTTP 405 Method Not Allowed
func NotAllowed(w http.ResponseWriter) {
	Error(w, http.StatusMethodNotAllowed, "Method Not Allowed")
}

func NotLoggedInError(w http.ResponseWriter) {
	Error(w, http.StatusForbidden, "You are not logged in")
}

// Error basic function to send any type of Error
func Error(w http.ResponseWriter, statusCode int, errMessage string) {
	// We need to set Content-Type here, cause we can not do it after WriteHeader
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(statusCode)

	jsonResponse(w, &ErrorResponse{
		ErrorData: ErrorData{
			Code:    statusCode,
			Message: errMessage,
		},
	})
}

func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", contentTypeJSON)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		panic(err)
	}
}
