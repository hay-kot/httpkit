// Package errtrace provides a simple way to annotate errors with the stack trace of the caller.
package errtrace

import (
	"errors"
	"fmt"
	"runtime"
)

func IsTraceable(err error) bool {
	t := &stacktrace{}
	return errors.As(err, &t)
}

// TraceWrap wraps an error within a stacktrace and returns the new error.
//   - The stacktrace is generated at the point of the call to TraceWrap.
//   - The message is formatted using fmt.Sprintf. The error returned
//     by TraceWrap implements the Error and Unwrap interfaces.
//   - If the error is nil, TraceWrap returns nil.
//
// TraceWrap errors are intended to be used for debugging and should only be
// viewed by developers. They should _generally_ not be returned to the users
// of your application.
//
// To use a TraceWrap error, you can use the TraceString function to get a
// printable string of the stacktrace.
//
// Example:
//
//	func doSomething(v string) error {
//	  err := doSomethingElse()
//	  if err != nil {
//	    return errtrace.TraceWrap(err, "failed to do something with %s", v)
//	  }
//
//	  return nil
//	}
func TraceWrap(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	return newTraceable(err, msg, args...)
}

// Trace creates a new error with a stacktrace and returns the new error.
// Use this like you would fmt.Errorf. For wrapping errors use TraceWrap.
func Trace(msg string, args ...any) error {
	return newTraceable(nil, msg, args...)
}

// MakeTraceable wraps an error within a stacktrace and returns the new error.
func MakeTraceable(err error) error {
	if err == nil {
		return nil
	}

	return newTraceable(err, err.Error())
}

func newTraceable(cause error, msg string, args ...any) error {
	err := &stacktrace{
		message: fmt.Sprintf(msg, args...),
		cause:   cause,
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return err
	}

	err.file = cleanGoPath(file)
	err.line = line

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return err
	}

	err.function = fn.Name()

	return err
}

// unwrapToTraceable unwraps the error and returns the first error that is a
// (*stacktrace). If no *stacktrace is found, nil is returned.
func unwrapToTraceable(err error) *stacktrace {
	for {
		if err == nil {
			return nil
		}

		traceable, ok := err.(*stacktrace) //nolint:errorlint
		if !ok {
			maybeErr := errors.Unwrap(err)
			if maybeErr == nil {
				return nil
			}
			err = maybeErr
			continue
		}

		return traceable
	}
}
