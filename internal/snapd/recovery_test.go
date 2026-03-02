package snapd_test

import (
	"context"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/matryer/is"
)

func TestCheckRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKey string

		want    bool
		wantErr bool
	}{
		"valid recovery key": {recoveryKey: "12345678-1234-5678-1234-567812345678", want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			

			c := snapd.New()

			// TODO: create the mock
			return

			got, err := c.CheckRecoveryKey(context.Background(), tc.recoveryKey, nil)
			if tc.wantErr {
				is.True(err != nil) // Expected an error but got nil
			}
			is.NoErr(err) // Unexpected error

			is.Equal(tc.want, got) // Got %v, want %v
		})
	}
}
