package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response implements the standard JSON response payload structure.
type Response struct {
	Status string          `json:"status"`
	Error  *ResponseError  `json:"error,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

// ResponseError implements the standard Error response structure.
type ResponseError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func (e ResponseError) Error() string {
	j, err := json.Marshal(e)
	if err != nil {
		return "ResponseError: " + err.Error()
	}
	return string(j)
}

// Fail ends an unsuccessful JSON response with the stardard failure format.
func Fail(w http.ResponseWriter, status, errCode int, details ...string) {
	msg, ok := frErrMap[errCode]
	if !ok {
		errCode = status
		msg = http.StatusText(status)
	}
	r := &Response{
		Status: StatusFail,
		Error: &ResponseError{
			Code:    errCode,
			Message: msg,
			Details: details,
		},
	}
	j, err := json.Marshal(r)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(j)
	if err != nil {
		log.Printf("couldn't write response: %v", err)
	}
}

// Send sends a successful JSON response using the standard success format.
func Send(w http.ResponseWriter, status int, result interface{}) {
	rj, err := json.Marshal(result)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	r := &Response{
		Status: StatusOK,
		Result: rj,
	}
	j, err := json.Marshal(r)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(j)
	if err != nil {
		log.Printf("couldn't write response: %v", err)
	}
}

// ResponseStatus constants.
const (
	StatusOK   = "ok"
	StatusFail = "nok"
)

// ErrCodes map response body error codes to standard error status messages.
const (
	ErrCodeAuto                   = 0
	ErrCodeParameterMismatch      = 1
	ErrCodeHMACInvalid            = 2
	ErrCodeOperationNotAllowed    = 3
	ErrCodeInvalidUserCredentials = 4
	ErrCodeInternalServiceError   = 666
	ErrCodeServiceUnavailable     = 667
)

var frErrMap = map[int]string{
	ErrCodeParameterMismatch:      "Parameter Mismatch",
	ErrCodeHMACInvalid:            "HMAC is invalid",
	ErrCodeOperationNotAllowed:    "Operation is not allowed",
	ErrCodeInvalidUserCredentials: "Invalid user credentials",
	ErrCodeInternalServiceError:   "Internal Server Error",
	ErrCodeServiceUnavailable:     "Service Unavailable",
}
