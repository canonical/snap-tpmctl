package cmd_test

import (
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
	"github.com/matryer/is"
)

func TestRemove(t *testing.T) {
	t.Parallel()

	commands := []string{
		"remove-passphrase",
		"remove-pin",
	}

	tests := map[string]struct {
		admineUID int

		wantErr bool
	}{
		"Success_on_removing": {},

		"Fail_on_user_privilege": {admineUID: 1, wantErr: true},
		"Fail_on_removing":       {wantErr: true},
	}
	for _, command := range commands {
		for name, tc := range tests {
			t.Run(filepath.Join(command, name), func(t *testing.T) {
				t.Parallel()

				is := is.New(t)
				ctx, logs := testutils.TestLoggerWithBuffer(t)

				var out strings.Builder
				tui := tui.New(nil, &out)

				c := snapdtestutils.NewMockSnapdServer(t, ctx)
				s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))
				app := cmd.New(
					cmdtestutils.WithSnapTPM(s),
					cmdtestutils.WithArgs(command),
					cmdtestutils.WithTui(tui),
					cmdtestutils.WithEuid(tc.admineUID),
				)

				err := app.Run(ctx)
				if testutils.CheckError(is, err, tc.wantErr) {
					return
				}

				is.True(logs.Len() == 0) // No logs printed by default

				golden.CheckOrUpdate(t, out.String()) // TestRemove returns the correct output
			})
		}
	}
}
