package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"snap-tpmctl/internal/log"
	"testing"

	"github.com/nalgeon/be"
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
		"Returns 1 when got an error": {app: mockApp{err: errors.New("desired error")}, want: 1, wantInLog: "desired error"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var logs bytes.Buffer
			out := io.MultiWriter(&logs, t.Output())
			ctx := log.WithContextLogger(context.Background(), out)

			got := run(ctx, tc.app)
			be.Equal(t, tc.want, got) // Return value does not match

			if tc.wantInLog != "" {
				if !bytes.Contains(logs.Bytes(), []byte(tc.wantInLog)) {
					t.Errorf("expected log to contain %q, but it didn't. Got: %s", tc.wantInLog, logs.String())
				}
			}
		})
	}
}
