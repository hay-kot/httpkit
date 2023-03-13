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
	const indent = "  "

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
		str.WriteRune('\n')
		str.WriteString(indent)
		str.WriteString(traceable.message)
		str.WriteRune('\n')

		str.WriteString(indent)
		str.WriteString("\033[4m")
		str.WriteString(traceable.file)
		str.WriteRune(':')
		str.WriteString(fmt.Sprintf("%d", traceable.line))
		str.WriteString("\033[0m")

		str.WriteString(indent)
		str.WriteString(traceable.function)
		str.WriteString("()")

		if traceable.cause != nil {
			str.WriteString(" -> ")
			str.WriteString(traceable.cause.Error())
			str.WriteRune('\n')
			err = traceable.cause
		}
	}

	return str.String()[:str.Len()-1] // remove last newline
}

func (st *stacktrace) String() string {
	return stringer(st)
}

type StackTraceData struct {
	Message  string
	File     string
	Function string
	Line     int
	Cause    error
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

func TraceString(err error) (trace string, ok bool) {
	if !IsTraceable(err) {
		return "", false
	}

	return stringer(err), true
}
