package cmd_test

import (
	"os"
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

func TestCreateKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKeyName string

		wantErr bool
	}{
		"Success_on_creting_recovery_key": {},

		"Error_from_snapd_on_empty_name":  {wantErr: true},
		"Error_from_snapd_on_unique_name": {recoveryKeyName: "test-duplicate", wantErr: true},
		"Error_on_creating_recovery_key":  {wantErr: true},
	}

	//nolint:dupl // regreneate and create have similar behaviour
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)
			ctx, logs := testutils.TestLoggerWithBuffer(t)

			command := "create-recovery-key"

			if tc.recoveryKeyName == "" {
				tc.recoveryKeyName = "test"
			}

			var out strings.Builder
			tui := tui.New(os.Stdin, &out)

			c := snapdtestutils.NewMockSnapdServer(t, ctx)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))
			app := cmd.New(
				cmdtestutils.WithSnapTPM(s),
				cmdtestutils.WithArgs(command, tc.recoveryKeyName),
				cmdtestutils.WithTui(tui),
			)

			err := app.Run(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default

			golden.CheckOrUpdate(t, out.String()) // TestCreateKey returns the expected output
		})
	}
}
