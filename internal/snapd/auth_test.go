package snapd_test

import (
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
	"github.com/canonical/snap-tpmctl/internal/snapd"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/matryer/is"
)

func TestReplacePassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		oldPassphrase string
		newPassphrase string

		wantErr bool
	}{
		"Passphrase_is_changed": {oldPassphrase: "old", newPassphrase: "new"},

		"Error_on_snapd_call_returning_400": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx, testutils.TestPath(t))

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

func TestCheckPassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		passphrase string

		wantErr bool
	}{
		"Passphrase_is_valid": {passphrase: "test12345"},

		"Error_on_low_quality_passphrase": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx, testutils.TestPath(t))

			err := c.CheckPassphrase(ctx, tc.passphrase)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
		})
	}
}

func TestCheckPIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pin string

		wantErr bool
	}{
		"PIN_is_valid": {pin: "12345"},

		"Error_on_low_quality_pin": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx, testutils.TestPath(t))

			err := c.CheckPIN(ctx, tc.pin)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
		})
	}
}

func TestReplacePIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		oldPIN string
		newPIN string

		wantErr bool
	}{
		"PIN_is_changed": {oldPIN: "12345", newPIN: "54321"},

		"Error_on_snapd_call_returning_400": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx, testutils.TestPath(t))

			// Container roles are passed as is to snapd. Not handled in that test.
			err := c.ReplacePIN(ctx, tc.oldPIN, tc.newPIN, nil)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
		})
	}
}

func TestReplacePlatformKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		authMode snapd.AuthMode
		secret   string

		wantErr bool
	}{
		"Replace_with_AuthModeNone":       {authMode: snapd.AuthModeNone},
		"Replace_with_AuthModePIN":        {authMode: snapd.AuthModePIN, secret: "12345"},
		"Replace_with_AuthModePassphrase": {authMode: snapd.AuthModePassphrase, secret: "test"},

		"Error_on_authModeNone_with_secret": {authMode: snapd.AuthModeNone, secret: "test", wantErr: true},
		"Error_on_snapd_call_returning_400": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			c := snapdtestutils.NewMockSnapdServer(t, ctx, testutils.TestPath(t))

			err := c.ReplacePlatformKey(ctx, tc.authMode, tc.secret)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
		})
	}
}
