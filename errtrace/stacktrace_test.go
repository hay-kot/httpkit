package errtrace

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func currentFileAbsolutePath() string {
	_, file, _, _ := runtime.Caller(0) // nolint:dogsled
	return file
}

func TestMarshalStack(t *testing.T) {
	filepath := currentFileAbsolutePath()
	funcname := "github.com/hay-kot/httpkit/errtrace.TestMarshalStack.func"

	tests := []struct {
		name  string
		error func() error
		want  []frame
	}{
		{
			name: "generic wrapped errors",
			error: func() error {
				err := errors.New("user with id 1 already exists")
				err = fmt.Errorf("wrap: %w", err)
				return err
			},
			want: []frame{
				{Error: "wrap: user with id 1 already exists"},
				{Error: "user with id 1 already exists"},
			},
		},
		{
			name: "all traceable errors",
			error: func() error {
				err := New("root error")
				err = Wrapf(err, "wrap 1")
				err = Wrapf(err, "wrap 2")

				return err
			},
			want: []frame{
				{Error: "wrap 2", Source: filepath, Line: 42, Function: funcname + "2"},
				{Error: "wrap 1", Source: filepath, Line: 41, Function: funcname + "2"},
				{Error: "root error", Source: filepath, Line: 40, Function: funcname + "2"},
			},
		},
		{
			name: "mixed traceable and non-traceable errors",
			error: func() error {
				err := New("user with id 1 already exists")
				err = fmt.Errorf("wrap: %w", err)
				err = Wrapf(err, "wrap 1")
				return err
			},
			want: []frame{
				{Error: "wrap 1", Source: filepath, Line: 57, Function: funcname + "3"},
				{Error: "wrap: user with id 1 already exists"},
				{Error: "user with id 1 already exists", Source: filepath, Line: 55, Function: funcname + "3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MarshalStack(tt.error())

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
