package tpm_test

import (
	"strings"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
	"github.com/canonical/snap-tpmctl/internal/snapd"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/matryer/is"
)

func TestAddPassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		passphrase string

		wantErr bool
	}{
		"Success_adding_passphrase": {passphrase: "test"},

		"Fail_adding_passphrase": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePlatformKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.AddPassphrase(ctx, tc.passphrase)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			var found bool
			for _, r := range *c.Requests {
				// check at least one request contains the expected passphrase
				if strings.Contains(r.Body, tc.passphrase) {
					found = true
					break
				}
			}
			is.True(found)
		})
	}

}

func TestReplacePassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		old string
		new string

		wantErr bool
	}{
		"Success_replacing_passphrase": {old: "old", new: "new"},

		"Fail_replacing_passphrase": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePassphrase")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.ReplacePassphrase(ctx, tc.old, tc.new)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			var found bool
			for _, r := range *c.Requests {
				// check at least one request contains the expected passphrases
				if strings.Contains(r.Body, tc.old) && strings.Contains(r.Body, tc.new) {
					found = true
					break
				}
			}
			is.True(found)
		})
	}
}

func TestRemovePassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Success_removing_passphrase": {},

		"Fail_removing_passphrase": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePlatformKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.RemovePassphrase(ctx)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			var found bool
			for _, r := range *c.Requests {
				// check at least one request contains the expected mode
				if strings.Contains(r.Body, string(snapd.AuthModeNone)) {
					found = true
					break
				}
			}
			is.True(found)
		})
	}

}

func TestAddPIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pin string

		wantErr bool
	}{
		"Success_adding_pin": {pin: "123456"},

		"Fail_adding_pin": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePlatformKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.AddPIN(ctx, tc.pin)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			var found bool
			for _, r := range *c.Requests {
				// check at least one request contains the expected pin
				if strings.Contains(r.Body, tc.pin) {
					found = true
					break
				}
			}
			is.True(found)
		})
	}
}

func TestReplacePIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		old string
		new string

		wantErr bool
	}{
		"Success_replacing_pin": {old: "123456", new: "654321"},

		"Fail_replacing_pin": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePIN")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.ReplacePIN(ctx, tc.old, tc.new)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			var found bool
			for _, r := range *c.Requests {
				// check at least one request contains the expected pins
				if strings.Contains(r.Body, tc.old) && strings.Contains(r.Body, tc.new) {
					found = true
					break
				}
			}
			is.True(found)
		})
	}
}

func TestRemovePIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Success_removing_pin": {},

		"Fail_removing_pin": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePlatformKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.RemovePIN(ctx)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			var found bool
			for _, r := range *c.Requests {
				// check at least one request contains the expected mode
				if strings.Contains(r.Body, string(snapd.AuthModeNone)) {
					found = true
					break
				}
			}
			is.True(found)
		})
	}
}

/*
func TestReplacePassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		oldPassphrase string
		newPassphrase string

		replacePassphraseError bool

		wantErr bool
	}{
		"Success": {oldPassphrase: "old-passphrase", newPassphrase: "new-passphrase"},

		"Error when snapd down": {oldPassphrase: "old-passphrase", newPassphrase: "new-passphrase", replacePassphraseError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePassphraseError: tc.replacePassphraseError,
			})

			err := tpm.ReplacePassphrase(ctx, mockClient, tc.oldPassphrase, tc.newPassphrase)

			if tc.wantErr {
				is.True(err != nil) // Expected an error but got nil
				return
			}
			is.NoErr(err) // Unexpected error
		})
	}
}

func TestReplacePIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		oldPIN string
		newPIN string

		replacePINError bool

		wantErr bool
	}{
		"Success": {oldPIN: "123456", newPIN: "654321"},

		"Error when snapd down": {oldPIN: "123456", newPIN: "654321", replacePINError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePINError: tc.replacePINError,
			})

			err := tpm.ReplacePIN(ctx, mockClient, tc.oldPIN, tc.newPIN)

			if tc.wantErr {
				is.True(err != nil) // Expected an error but got nil
				return
			}
			is.NoErr(err) // Unexpected error
		})
	}
}

func TestAddPIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		replacePlatformKeyError bool

		wantErr bool
	}{
		"Adds PIN authentication": {},

		"Error when snapd down": {replacePlatformKeyError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
			})

			err := tpm.AddPIN(ctx, mockClient, "123456")

			if tc.wantErr {
				is.True(err != nil) // Expected an error but got nil
				return
			}
			is.NoErr(err) // Unexpected error
		})
	}
}

func TestRemovePIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		replacePlatformKeyError bool
		replacePlatformKeyNotOK bool

		wantErr bool
	}{
		"Removes PIN authentication": {},

		"Error when snapd down": {replacePlatformKeyError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
			})

			err := tpm.RemovePIN(ctx, mockClient)

			if tc.wantErr {
				is.True(err != nil) // Expected an error but got nil
				return
			}
			is.NoErr(err) // Unexpected error
		})
	}
}

func TestAddPassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		replacePlatformKeyError bool
		replacePlatformKeyNotOK bool

		wantErr bool
	}{
		"Adds passphrase authentication": {},

		"Error when snapd down": {replacePlatformKeyError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
			})

			err := tpm.AddPassphrase(ctx, mockClient, "my-secure-passphrase")

			if tc.wantErr {
				is.True(err != nil) // Expected an error but got nil
				return
			}
			is.NoErr(err) // Unexpected error
		})
	}
}

func TestRemovePassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		replacePlatformKeyError bool
		replacePlatformKeyNotOK bool

		wantErr bool
	}{
		"Removes passphrase authentication": {},

		"Error when snapd down": {replacePlatformKeyError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
			})

			err := tpm.RemovePassphrase(ctx, mockClient)

			if tc.wantErr {
				is.True(err != nil) // Expected an error but got nil
				return
			}
			is.NoErr(err) // Unexpected error
		})
	}
}
*/
