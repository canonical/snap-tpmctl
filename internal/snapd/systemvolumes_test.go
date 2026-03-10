package snapd_test

import (
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
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
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx, testutils.TestPath(t))

			got, err := c.ListVolumeInfo(ctx)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			golden.CheckOrUpdateYAML(t, got)
		})
	}
}
