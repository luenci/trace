package trace_test

import (
	trace "github.com/luenci/trace"
	"testing"
)

func TestTrace(t *testing.T) {
	tests := []struct {
		name string
		want func()
	}{
		{name: "tests", want: func() {
			trace.Trace()
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want()
		})
	}
}
