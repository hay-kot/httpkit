package server

import (
	"encoding/json"
	"net/http"
)

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(val); err != nil {
		return err
	}
	return nil
}

// DecodeStrict reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
// DecodeStrict will return an error if the JSON document contains fields that
// are not defined in the provided value.
func DecodeStrict(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}
	return nil
}
