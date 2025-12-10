package tpm_test

import (
	"context"
	"testing"

	"github.com/nalgeon/be"
	"snap-tpmctl/internal/testutils"
	"snap-tpmctl/internal/tpm"
)

//nolint:dupl // Similar test pattern for different function (ReplacePassphrase vs ReplacePIN)
func TestReplacePassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		oldPassphrase string
		newPassphrase string

		replacePassphraseError bool
		replacePassphraseNotOK bool

		wantErr bool
	}{
		"Success": {oldPassphrase: "old-passphrase", newPassphrase: "new-passphrase"},

		"Error when snapd down":      {oldPassphrase: "old-passphrase", newPassphrase: "new-passphrase", replacePassphraseError: true, wantErr: true},
		"Error when response not ok": {oldPassphrase: "old-passphrase", newPassphrase: "new-passphrase", replacePassphraseNotOK: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePassphraseError: tc.replacePassphraseError,
				ReplacePassphraseNotOK: tc.replacePassphraseNotOK,
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

//nolint:dupl // Similar test pattern for different function (ReplacePassphrase vs ReplacePIN)
func TestReplacePIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		oldPin string
		newPin string

		replacePINError bool
		replacePINNotOK bool

		wantErr bool
	}{
		"Success": {oldPin: "123456", newPin: "654321"},

		"Error when snapd down":      {oldPin: "123456", newPin: "654321", replacePINError: true, wantErr: true},
		"Error when response not ok": {oldPin: "123456", newPin: "654321", replacePINNotOK: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePINError: tc.replacePINError,
				ReplacePINNotOK: tc.replacePINNotOK,
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
		replacePlatformKeyNotOK bool

		wantErr bool
	}{
		"Adds PIN authentication": {},

		"Error when snapd down":      {replacePlatformKeyError: true, wantErr: true},
		"Error when response not ok": {replacePlatformKeyNotOK: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
				ReplacePlatformKeyNotOK: tc.replacePlatformKeyNotOK,
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

		"Error when snapd down":      {replacePlatformKeyError: true, wantErr: true},
		"Error when response not ok": {replacePlatformKeyNotOK: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
				ReplacePlatformKeyNotOK: tc.replacePlatformKeyNotOK,
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

		"Error when snapd down":      {replacePlatformKeyError: true, wantErr: true},
		"Error when response not ok": {replacePlatformKeyNotOK: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
				ReplacePlatformKeyNotOK: tc.replacePlatformKeyNotOK,
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

		"Error when snapd down":      {replacePlatformKeyError: true, wantErr: true},
		"Error when response not ok": {replacePlatformKeyNotOK: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				ReplacePlatformKeyError: tc.replacePlatformKeyError,
				ReplacePlatformKeyNotOK: tc.replacePlatformKeyNotOK,
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
