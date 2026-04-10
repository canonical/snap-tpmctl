package cmd_test

import (
	"fmt"
	"path/filepath"
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

func TestAdd(t *testing.T) {
	t.Parallel()

	commands := []string{
		"add-passphrase",
		"add-pin",
	}

	tests := map[string]struct {
		wantUserErr bool
		wantReadErr bool
		wantErr     bool
	}{
		"Success_on_adding": {},

		"Fail_on_user_privilege": {wantUserErr: true, wantErr: true},
		"Fail_reading_input":     {wantReadErr: true, wantErr: true},
		"Fail_wrong_auth_mode":   {wantErr: true},
		"Fail_on_validating":     {wantErr: true},
		"Fail_on_adding":         {wantErr: true},
	}
	for _, command := range commands {
		for name, tc := range tests {
			t.Run(filepath.Join(command, name), func(t *testing.T) {
				t.Parallel()

				is := is.New(t)
				ctx, logs := testutils.TestLoggerWithBuffer(t)

				input := "test"
				if command == "add-pin" {
					input = "12345"
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
					for range 2 {
						fmt.Fprintf(ptmx, "%s\n", input)
					}
				}()

				euid := 0
				if tc.wantUserErr {
					euid = 1
				}

				c := snapdtestutils.NewMockSnapdServer(t, ctx)
				s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))
				app := cmd.New(
					cmdtestutils.WithSnapTPM(s),
					cmdtestutils.WithArgs(command),
					cmdtestutils.WithTui(tui),
					cmdtestutils.WithEuid(euid),
				)

				err = app.Run(ctx)
				if testutils.CheckError(is, err, tc.wantErr) {
					return
				}

				is.True(logs.Len() == 0) // No logs printed by default

				golden.CheckOrUpdate(t, out.String()) // TestAdd retruns the correct output
			})
		}
	}
}
