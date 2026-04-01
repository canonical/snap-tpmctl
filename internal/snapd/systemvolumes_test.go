package snapd_test

import (
	"testing"

	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/matryer/is"
)

func TestListVolumeInfo(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Returns_volumes_info": {},
		"No_volumes":           {},

		"Error_on_invalid_result":           {wantErr: true},
		"Error_on_snapd_call_returning_400": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			c := snapdtestutils.NewMockSnapdServer(t, ctx)

			got, err := c.ListVolumeInfo(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			golden.CheckOrUpdate(t, got)
		})
	}
}
