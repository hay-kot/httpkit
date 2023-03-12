package server

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// JSON converts a Go value to JSON and sends it to the client.
// If the code is StatusNoContent, no body is sent.
// Adapted from https://github.com/ardanlabs/service/tree/master/foundation/web
func JSON(w http.ResponseWriter, code int, data interface{}) error {
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
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
