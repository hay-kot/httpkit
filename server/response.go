package server

import (
	"encoding/json"
	"net/http"
)

// MarshalFunc encodes a Go value to JSON.
type MarshalFunc func(v any) ([]byte, error)

var marshal MarshalFunc = json.Marshal

// SetJSONMarshaler sets the function JSON and ErrorBuilder.Write use to encode
// response bodies. The default is encoding/json.Marshal. Use it to swap in another
// encoder, such as encoding/json/v2 with custom options. It is not safe to call
// concurrently with request handling, so set it once at startup.
func SetJSONMarshaler(fn MarshalFunc) {
	marshal = fn
}

// JSON converts a Go value to JSON and sends it to the client.
// If the code is StatusNoContent, no body is sent.
// Adapted from https://github.com/ardanlabs/service/tree/master/foundation/web
func JSON(w http.ResponseWriter, code int, data interface{}) error {
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return nil
	}

	jsonData, err := marshal(data)
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
