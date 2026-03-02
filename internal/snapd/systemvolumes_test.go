package snapd_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
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
		"No_volume_info": {},
		
		"Error_on_snapd_call_returning_400": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			resp, err := os.ReadFile(testutils.TestPath(t))
			is.NoErr(err) // Setup: read the test reponse from test file asset

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v2/system-volumes" {
					http.NotFound(w, r)
					return
				}
				_, err := w.Write (resp)
				is.NoErr(err) // Server: could not write response to client
			}))
			defer ts.Close()


			c := snapd.New(snapdtestutils.WithBaseURL(ts.URL))


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
