package snapd_test

import (
	"context"
	"testing"

	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/matryer/is"
)

func TestGenerateRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Returns_recovery_key_result": {},

		"Error_on_invalid_result":           {wantErr: true},
		"Error_on_snapd_call_returning_400": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			c := snapdtestutils.NewMockSnapdServer(t, "/v2/system-volumes")

			got, err := c.GenerateRecoveryKey(context.Background())
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			golden.CheckOrUpdateYAML(t, got)
		})
	}
}

func TestCheckRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKey string

		want    bool
		wantErr bool
	}{
		"Recovery_key_matches":                 {recoveryKey: "12345678-1234-5678-1234-567812345678", want: true},
		"Recovery_key_does_not_match":          {recoveryKey: "99999999-1234-5678-1234-567812345678", want: false},
		"Return_false_on_invalid_recovery_key": {recoveryKey: "invalid-format", want: false},

		"Error_on_invalid_input": {recoveryKey: "", wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			c := snapdtestutils.NewMockSnapdServer(t, "/v2/system-volumes")

			// Container roles are passed as is to snapd. Not handled in that test.
			valid, err := c.CheckRecoveryKey(context.Background(), tc.recoveryKey, nil)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.Equal(valid, tc.want)
		})
	}
}
