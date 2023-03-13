package errtrace

import (
	"errors"
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
		wantStr: "annotate_test.go:34 TestAnnotate_ErrorArg: prefix: wrapped error text",
	}

	err := Annotate(test.args.s, test.args.args...)
	if (err != nil) != test.wantErr {
		t.Errorf("Annotate() error = %v, wantErr %v", err, test.wantErr)
		return
	}

	if err != nil && err.Error() != test.wantStr {
		t.Errorf("Annotate() error = %v, want %v", err.Error(), test.wantStr)
	}
}
