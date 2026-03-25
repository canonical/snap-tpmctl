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

func TestStatus(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Returns_FDE_status": {},

		"Error_when_getting_FDE_status": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)
			ctx, logs := testutils.TestLoggerWithBuffer(t)

			command := "status"

			c := snapdtestutils.NewMockSnapdServer(t, ctx)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))
			app := cmd.New(cmdtestutils.WithSnapTPM(s), cmdtestutils.WithArgs(command))

			err := app.Run(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default
		})
	}
}
