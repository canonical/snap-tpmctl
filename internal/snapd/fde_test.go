package snapd_test

import (
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
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
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx, testutils.TestPath(t))

			got, err := c.FdeStatus(ctx)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.Equal(got, tc.want) // TestFDEStatus returns the expected FDE status
		})
	}
}
