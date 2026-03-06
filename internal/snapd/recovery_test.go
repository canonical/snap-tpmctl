package snapd_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
	"github.com/canonical/snap-tpmctl/internal/snapd"
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
			ctx := log.WithLoggerInContext(context.Background(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx)

			got, err := c.GenerateRecoveryKey(ctx)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			golden.CheckOrUpdateYAML(t, got)
		})
	}
}

func TestAddRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		keyID           string
		recoveryKeyName string

		wantErr bool
	}{
		"Returns_accepted": {keyID: "OVJe6EHITg", recoveryKeyName: "test"},

		"Error_on_invalid_key_id":            {keyID: "invalid-key-id", recoveryKeyName: "test", wantErr: true},
		"Error_on_invalid_recovery_key_name": {keyID: "OVJe6EHITg", recoveryKeyName: "default", wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			// TODO: generalize that
			ctx := log.WithLoggerInContext(context.Background(), t.Output())
			log.SetLoggerLevelInContext(ctx, slog.LevelDebug)

			c := snapdtestutils.NewMockSnapdServer(t, ctx)

			keySlots := []snapd.Keyslot{{Name: tc.recoveryKeyName}}
			err := c.AddRecoveryKey(ctx, tc.keyID, keySlots)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
		})
	}
}

func TestReplaceRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		keyID           string
		recoveryKeyName string

		wantErr bool
	}{
		"Returns_accepted": {keyID: "OVJe6EHITg", recoveryKeyName: "test"},

		"Error_on_invalid_key_id":            {keyID: "invalid-key-id", recoveryKeyName: "test", wantErr: true},
		"Error_on_invalid_recovery_key_name": {keyID: "OVJe6EHITg", recoveryKeyName: "default", wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			// TODO: generalize that
			ctx := log.WithLoggerInContext(context.Background(), t.Output())
			log.SetLoggerLevelInContext(ctx, slog.LevelDebug)

			c := snapdtestutils.NewMockSnapdServer(t, ctx)

			keySlots := []snapd.Keyslot{{Name: tc.recoveryKeyName}}
			err := c.ReplaceRecoveryKey(ctx, tc.keyID, keySlots)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
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

			ctx := log.WithLoggerInContext(context.Background(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx)

			// Container roles are passed as is to snapd. Not handled in that test.
			valid, err := c.CheckRecoveryKey(ctx, tc.recoveryKey, nil)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.Equal(valid, tc.want)
		})
	}
}
