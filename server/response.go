package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// JSON converts a Go value to JSON and sends it to the client.
// If the code is StatusNoContent, no body is sent.
// Adapted from https://github.com/ardanlabs/service/tree/master/foundation/web
func JSON(w http.ResponseWriter, code int, data interface{}) error {
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Write the status code to the response.
	w.WriteHeader(code)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

type responseError struct {
	Err string
}

func (e responseError) Error() string {
	return e.Err
}

// IsResponseError returns true if the error is a ResponseError or wraps a ResponseError. If
// the error is a ResponseError, the caller can assume that the HTTP status code and message
// have already been sent to the client.
func IsResponseError(err error) bool {
	if err == nil {
		return false
	}

	var e responseError
	return errors.As(err, &e)
}

type errorResp struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Error      string            `json:"error"`
	Context    map[string]string `json:"context,omitempty"`
}

type ErrorOption func(*errorResp)

// WithContext adds a context map to the error response.
//
// Example JSON:
//
//	{
//	  "error": "invalid request body",
//	  "context": {
//	    "field": "body",
//	    "reason": "missing required field"
//	  }
//	}
func WithContext(ctx map[string]string) ErrorOption {
	return func(e *errorResp) {
		e.Context = ctx
	}
}

// Error sends an error response back to the client with the status code and
// message provided.
//
// Example:
//
//	return server.Error(w, http.StatusBadRequest, "invalid request body")
//
// ExampleJSON:
//
//		{
//		 "error": "invalid request body"
//	   "statusCode": 400
//		}
func Error(w http.ResponseWriter, code int, msg string, fns ...ErrorOption) error {
	body := errorResp{Error: msg, StatusCode: code}
	for _, fn := range fns {
		fn(&body)
	}

	err := JSON(w, code, body)
	if err != nil {
		return err
	}

	return responseError{Err: msg}
}

// Errorf is a convenience function that formats the message and calls Error.
func Errorf(w http.ResponseWriter, code int, msg string, args ...any) error {
	msg = fmt.Sprintf(msg, args...)
	return Error(w, code, msg)
}
