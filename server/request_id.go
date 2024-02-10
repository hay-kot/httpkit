package server

import "context"

type RequestIDFunc func(ctx context.Context) string

var requestIDFunc RequestIDFunc = func(ctx context.Context) string {
	return ""
}

// SetRequestIDFunc sets the function used to get the request ID from a context.
func SetRequestIDFunc(fn RequestIDFunc) {
	requestIDFunc = fn
}
