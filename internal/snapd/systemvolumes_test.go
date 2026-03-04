package snapd_test

import (
	"context"
	"testing"

	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/matryer/is"
)

/*

{
  "result": {
    "message": "this action is not supported on this system"
  },
  "status": "Bad Request",
  "status-code": 400,
  "type": "error"
}

*/

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
			c := snapdtestutils.NewMockSnapdServer(t, "/v2/system-volumes")

			got, err := c.ListVolumeInfo(context.Background())
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			golden.CheckOrUpdateYAML(t, got)
		})
	}
}
