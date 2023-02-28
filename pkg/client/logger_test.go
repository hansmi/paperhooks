package client

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type appendLogger []string

var _ stdLogger = (*appendLogger)(nil)

func (l *appendLogger) Print(v ...any) {
	*l = append(*l, fmt.Sprint(v...))
}

func TestLogger(t *testing.T) {
	for _, tc := range []struct {
		name  string
		build func(stdLogger) Logger
		want  []string
	}{
		{
			name: "discard",
			build: func(stdLogger) Logger {
				return &discardLogger{}
			},
		},
		{
			name: "wrapped",
			build: func(s stdLogger) Logger {
				return &wrappedStdLogger{s}
			},
			want: []string{
				"[E] error 1",
				"[W] warn 2",
				"[D] debug 3",
			},
		},
		{
			name: "prefix",
			build: func(s stdLogger) Logger {
				return &prefixLogger{
					wrapped: &wrappedStdLogger{s},
					prefix:  "test: ",
				}
			},
			want: []string{
				"[E] test: error 1",
				"[W] test: warn 2",
				"[D] test: debug 3",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var dest appendLogger

			logger := tc.build(&dest)
			logger.Errorf("error %d", 1)
			logger.Warnf("warn %d", 2)
			logger.Debugf("debug %d", 3)

			got := []string(dest)

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Log message diff (-want +got):\n%s", diff)
			}
		})
	}
}
