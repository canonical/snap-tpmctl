package cmd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	cmdtestutils "github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd/testutils"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/creack/pty"
	"github.com/matryer/is"
)

func TestCheck(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		key string

		wantReadErr bool
		wantErr     bool
	}{
		"Success_checking_recovery_key":       {},
		"Success_on_not_working_recovery_key": {},

		"Fail_reading_input":             {wantReadErr: true, wantErr: true},
		"Fail_with_invalid_recovery_key": {key: "invalid", wantErr: true},
		"Fail_checking_recovery_key":     {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)
			ctx, logs := testutils.TestLoggerWithBuffer(t)

			command := "check-recovery-key"
			if tc.key == "" {
				tc.key = "11272-47509-28031-54818-41671-38673-11053-06376"
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
				fmt.Fprintf(ptmx, "%s\n", tc.key)
			}()

			c := snapdtestutils.NewMockSnapdServer(t, ctx)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))
			app := cmd.New(
				cmdtestutils.WithSnapTPM(s),
				cmdtestutils.WithArgs(command),
				cmdtestutils.WithTui(tui),
			)

			err = app.Run(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default

			golden.CheckOrUpdate(t, out.String()) // TestCheck returns the expected output
		})
	}
}
