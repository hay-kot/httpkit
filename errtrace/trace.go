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

// Trace creates a new error with a stacktrace and returns the new error.
// Use this like you would fmt.Errorf.
//
// Deprecated: Use New instead.
func Trace(msg string, args ...any) error {
	return newTraceable(nil, msg, args...)
}

// New creates a new error with a stacktrace and returns the new error.
// Use this like you would fmt.Errorf.
func New(msg string, args ...any) error {
	return newTraceable(nil, msg, args...)
}

// TraceWrap is the same as Trace, but it wraps an existing error.
// Deprecated: Use Wrapf instead.
func TraceWrap(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	return newTraceable(err, msg, args...)
}

// Wrap wraps an error within a stacktrace and returns the new error.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	return newTraceable(err, err.Error())
}

// Wrapf wraps an error within a stacktrace and returns the new error.
//   - The stacktrace is generated at the point of the call to Wrapf.
//   - The message is formatted using fmt.Sprintf. The error returned
//     by Wrapf implements the Error and Unwrap interfaces.
//   - If the error is nil, Wrapf returns nil.
//
// Wrapf errors are intended to be used for debugging and should only be
// viewed by developers. They should _generally_ not be returned to the users
// of your application.
//
// To use a Wrapf error, you can use the TraceString function to get a
// printable string of the stacktrace.
//
// Example:
//
//	func doSomething(v string) error {
//	  err := doSomethingElse()
//	  if err != nil {
//	    return errtrace.Wrapf(err, "failed to do something with %s", v)
//	  }
//
//	  return nil
//	}
func Wrapf(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	return newTraceable(err, msg, args...)
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
