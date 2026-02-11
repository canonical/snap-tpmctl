package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"

	"github.com/matryer/is"
)

type mockApp struct{ err error }

func (m mockApp) Run(ctx context.Context) error {
	return m.err
}

func TestRun(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		app mockApp

		want      int
		wantInLog string
	}{
		"Returns 0 on success":        {app: mockApp{err: nil}, want: 0},
		"Returns 1 when got an error": {app: mockApp{err: errors.New("desired error")}, want: 1, wantInLog: "desired errors"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.NewRelaxed(t)

			var logs bytes.Buffer
			w := io.MultiWriter(&logs, t.Output())
			ctx := log.WithLoggerInContext(context.Background(), w)

			got := run(ctx, tc.app)
			is.Equal(tc.want, got) // Return value does not match exit code

			if tc.wantInLog != "" {
				is.True(strings.Contains(logs.String(), tc.wantInLog)) // Log does not contain expected message
			}
		})
	}
}
