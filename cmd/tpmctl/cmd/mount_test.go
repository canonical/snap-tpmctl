package cmd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	cmdtestutils "github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/creack/pty"
	"github.com/matryer/is"
)

func TestGetLuksKeyFromRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKey string
		hexFlag     bool
		escapedFlag bool

		wantReadErr bool
		wantErr     bool
	}{
		"Success_getting_luks_key":                   {},
		"Success_getting_luks_key_with_hex_flag":     {hexFlag: true},
		"Success_getting_luks_key_with_escaped_flag": {escapedFlag: true},

		"Fail_reading_input":      {wantReadErr: true, wantErr: true},
		"Fail_getting_luks_key":   {recoveryKey: "invalid", wantErr: true},
		"Fail_for_too_many_flags": {wantErr: true, hexFlag: true, escapedFlag: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)
			ctx, logs := testutils.TestLoggerWithBuffer(t)

			command := "get-luks-key"

			args := []string{command}
			if tc.hexFlag {
				args = append(args, "--hex")
			}
			if tc.escapedFlag {
				args = append(args, "--escaped")
			}

			if tc.recoveryKey == "" {
				tc.recoveryKey = "11272-47509-28031-54818-41671-38673-11053-06376"
			}

			ptmx, tty, err := pty.Open()
			is.NoErr(err)
			defer ptmx.Close()
			defer tty.Close()

			if tc.wantReadErr {
				tty = nil
			}

			var out strings.Builder
			tui := tui.New(tty, &out)

			done := make(chan struct{})
			go func() {
				defer close(done)
				fmt.Fprintf(ptmx, "%s\n", tc.recoveryKey)
			}()

			app := cmd.New(
				cmdtestutils.WithArgs(args...),
				cmdtestutils.WithTui(tui),
			)

			err = app.Run(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default

			golden.CheckOrUpdate(t, out.String()) // TestGetLuksKeyFromRecoveryKey returns the expected output
		})
	}
}
