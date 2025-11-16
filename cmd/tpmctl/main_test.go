package main

import (
	"errors"
	"testing"
)

type mockApp struct{ err error }

func (m mockApp) Run() error { return m.err }

func TestRun_Table(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		app            mockApp
		wantExit       int
		wantLogContain string
	}{
		"exit_with_success":      {app: mockApp{err: nil}, wantExit: 0},
		"exit_with_error_code_1": {app: mockApp{err: errors.New("foo")}, wantExit: 1, wantLogContain: "foo"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tc.wantExit != 0 {
			}

			if got := run(tc.app); got != tc.wantExit {
				t.Fatalf("run() = %d, want %d", got, tc.wantExit)
			}
		})
	}
}
