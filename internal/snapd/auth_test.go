package snapd_test

import (
	"context"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/matryer/is"
)

func TestReplacePassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		oldPassphrase string
		newPassphrase string

		wantErr bool
	}{
		"Passphrase_is_changed": {oldPassphrase: "test", newPassphrase: "test2"},

		"Error_on_snapd_call_returning_400": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(context.Background(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx)

			// Container roles are passed as is to snapd. Not handled in that test.
			err := c.ReplacePassphrase(ctx, tc.oldPassphrase, tc.newPassphrase, nil)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
		})
	}
}
