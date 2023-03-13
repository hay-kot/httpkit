package errtrace

import (
	"testing"
)

func Test_trimFuncName(t *testing.T) {
	tests := []struct {
		name string
		fn   string
		want string
	}{
		{
			name: "empty",
			fn:   "",
			want: "",
		},
		{
			name: "pointer receiver",
			fn:   "github.com/alexflint/go-arg/pkg.(*Receiver).FuncName",
			want: "(*Receiver).FuncName",
		},
		{
			name: "value receiver",
			fn:   "github.com/test/package/pkg.Receiver.FuncName",
			want: "Receiver.FuncName",
		},
		{
			name: "package func",
			fn:   "github.com/test/package/pkg.FuncName",
			want: "FuncName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimFuncName(tt.fn); got != tt.want {
				t.Errorf("trimFuncName() = %v, want %v", got, tt.want)
			}
		})
	}
}
