package errtrace

import (
	"errors"
	"regexp"
	"testing"
)

type annotationTest struct {
	name string
	args struct {
		s    string
		args []any
	}
	wantErr bool
	wantStr string
}

func TestAnnotate_ErrorArg(t *testing.T) {
	test := annotationTest{
		name: "error arg",
		args: struct {
			s    string
			args []any
		}{
			s: "prefix: %w",
			args: []any{
				errors.New("wrapped error text"),
			},
		},
		wantErr: true,
		wantStr: ".*annotate_test.go:35 TestAnnotate_ErrorArg: prefix: wrapped error text",
	}

	err := Annotatef(test.args.s, test.args.args...)
	if (err != nil) != test.wantErr {
		t.Errorf("Annotate() error = %v, wantErr %v", err, test.wantErr)
		return
	}

	regexp := regexp.MustCompile(test.wantStr)
	if !regexp.MatchString(err.Error()) {
		t.Errorf("Annotate() error = %v, wantErr %v", err, test.wantStr)
		return
	}
}
