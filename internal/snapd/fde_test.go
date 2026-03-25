package snapd_test

import (
	"testing"

	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/matryer/is"
)

func TestFdeStatus(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		want    string
		wantErr bool
	}{
		"Returns_FDE_status": {want: "enabled"},

		"Error_on_invalid_result":       {wantErr: true},
		"Error_when_getting_FDE_status": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			c := snapdtestutils.NewMockSnapdServer(t, ctx)

			got, err := c.FdeStatus(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.Equal(got, tc.want) // TestFDEStatus returns the expected FDE status
		})
	}
}
