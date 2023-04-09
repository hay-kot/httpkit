package errtrace

import (
	"errors"
	"fmt"
	"strings"
)

type stacktrace struct {
	message  string
	cause    error
	file     string
	function string
	line     int
}

func (st *stacktrace) Error() string {
	return st.message
}

func (st *stacktrace) Unwrap() error {
	return st.cause
}

func red(s string) string  { return "\033[31m" + s + "\033[0m" }
func bold(s string) string { return "\033[1m" + s + "\033[0m" }

func stringer(err error) string {
	str := strings.Builder{}
	lastWasTraceable := false
	const indent = "    "

	for {
		if err == nil {
			break
		}

		traceable, ok := err.(*stacktrace) //nolint:errorlint
		if !ok {
			if !lastWasTraceable {
				str.WriteString(red(bold("error: ")))
				str.WriteRune('\n')
				str.WriteString(indent)
				str.WriteString(err.Error())
				str.WriteRune('\n')
			}

			maybeErr := errors.Unwrap(err)
			if maybeErr == nil {
				break
			}

			lastWasTraceable = false
			err = maybeErr

			continue
		}

		lastWasTraceable = true

		str.WriteString(red(bold(("trace error: "))))

		str.WriteString(traceable.message)
		str.WriteRune('\n')

		// File and line numbers
		str.WriteString(indent)
		str.WriteString(traceable.file)
		str.WriteRune(':')
		str.WriteString(fmt.Sprintf("%d", traceable.line))

		// Function
		str.WriteRune('\n')
		str.WriteString(indent)
		str.WriteString(indent)
		str.WriteString(traceable.function)
		str.WriteString("()")

		if traceable.cause != nil {
			str.WriteString(" -> ")
			str.WriteString(traceable.cause.Error())
			str.WriteRune('\n')
		}

		err = traceable.cause
	}

	return str.String()
}

func (st *stacktrace) String() string {
	return stringer(st)
}

// StackTraceData contains the data of a stacktrace. It is returned by TraceData.
// It is not meant to be used directly, but rather to be used by other packages.
type StackTraceData struct {
	Message  string // The message of the error
	File     string // The file of the caller of the Traceable function
	Function string // The function of the caller of the Traceable function
	Line     int    // The line number of the caller of the Traceable function
	Cause    error  // The underlying error wrapped by the Traceable function
}

// Loc returns a formatted string that contains the file, function and line number of the caller of the Traceable function.
//
// Example: "main.go:42 main.main"
func (s *StackTraceData) Loc() string {
	return fmt.Sprintf("%s:%d %s", s.File, s.Line, s.Function)
}

// TraceData returns the file, function and line number of the caller of the Traceable function.
func TraceData(err error) (*StackTraceData, error) {
	if !IsTraceable(err) {
		return nil, errors.New("error is not traceable")
	}

	trace := unwrapToTraceable(err)
	if trace == nil {
		return nil, errors.New("error is not traceable")
	}

	return &StackTraceData{
		Message:  trace.message,
		File:     trace.file,
		Function: trace.function,
		Line:     trace.line,
		Cause:    trace.cause,
	}, nil
}

// TraceString returns a string representation of the error and all its causes.
func TraceString(err error) string {
	return stringer(err)
}
