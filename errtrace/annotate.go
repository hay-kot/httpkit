package errtrace

import (
	"errors"
	"fmt"
	"runtime"
)

// Annotate is a lighter version of Traceable and a drop in replacement for fmt.Errorf.
// instead of returning a traceable error it returns a standard error which is prefixed
// with the file and line number of the caller of Annotate.
//
// This is useful for annotating errors with additional file/line context, but without
// the overhead of a traceable error. Using Annotate is more appropriate for errors created
// by 'library' type code within your project.
//
// Example:
//
//	err := Annotate("users prefix: %w", errors.New("wrapped error text"))
//	println(err.Error())
//	// output: relative/path/filename.go:41 : prefix: wrapped error text
//
// Note:
//   - If the first argument is an empty string Annotate will return nil.
//   - If the args slice is empty Annotate will return nil.
func Annotate(s string, args ...any) error {
	if s == "" || len(args) == 0 {
		return nil
	}

	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf(s, args...)
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Errorf(s, args...)
	}

	return fmt.Errorf("%s:%d %s: "+s, append([]any{cleanGoPath(file), line, trimFuncName(fn.Name())}, args...)...)
}

// Annotated is a lighter version of Traceable and a drop in replacement for errors.New.
// instead of returning a traceable error it returns a standard error which is prefixed
// with the file and line number of the caller of Annotated.
func Annotated(s string) error {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return errors.New(s)
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return errors.New(s)
	}

	return fmt.Errorf("%s:%d %s: "+s, cleanGoPath(file), line, trimFuncName(fn.Name()))
}
