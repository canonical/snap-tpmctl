package tpm_test

import (
	"context"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/nalgeon/be"
)

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

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePassphraseError: tc.replacePassphraseError,
			})

			err := tpm.ReplacePassphrase(ctx, mockClient, tc.oldPassphrase, tc.newPassphrase)

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
		})
	}
}

func TestReplacePIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		oldPin string
		newPin string

		replacePINError bool

		wantErr bool
	}{
		"Success": {oldPin: "123456", newPin: "654321"},

		"Error when snapd down": {oldPin: "123456", newPin: "654321", replacePINError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePINError: tc.replacePINError,
			})

			err := tpm.ReplacePIN(ctx, mockClient, tc.oldPin, tc.newPin)

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
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

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
			})

			err := tpm.AddPIN(ctx, mockClient, "123456")

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
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

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
			})

			err := tpm.RemovePIN(ctx, mockClient)

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
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

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
			})

			err := tpm.AddPassphrase(ctx, mockClient, "my-secure-passphrase")

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
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

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
			})

			err := tpm.RemovePassphrase(ctx, mockClient)

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
		})
	}
}
