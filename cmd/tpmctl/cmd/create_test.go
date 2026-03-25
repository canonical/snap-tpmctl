package cmd_test

import (
	"testing"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	cmdtestutils "github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd/testutils"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/matryer/is"
)

func TestCreateKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKeyName string

		wantErr bool
	}{
		"Success_on_creting_recovery_key": {recoveryKeyName: "test"},

		"Error_on_empty_name":            {wantErr: true},
		"Error_on_unique_name":           {recoveryKeyName: "test-duplicate", wantErr: true},
		"Error_on_creating_recovery_key": {recoveryKeyName: "test", wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)
			ctx, logs := testutils.TestLoggerWithBuffer(t)

			command := "create-recovery-key"

			c := snapdtestutils.NewMockSnapdServer(t, ctx)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))
			app := cmd.New(cmdtestutils.WithSnapTPM(s), cmdtestutils.WithArgs(command, tc.recoveryKeyName))

			err := app.Run(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default
		})
	}
}
