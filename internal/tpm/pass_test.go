package tpm_test

import (
	"context"
	"testing"

	"github.com/nalgeon/be"
	"snap-tpmctl/internal/testutils"
	"snap-tpmctl/internal/tpm"
)

func TestIsValidPassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		passphrase string
		confirm    string

		checkPassphraseError   bool
		passphraseLowEntropy   bool
		passphraseInvalid      bool
		passphraseUnsupported  bool
		passphraseUnknownError bool
		passphraseNotOK        bool

		wantErr bool
	}{
		"Success": {},

		"Error when passphrase empty":           {wantErr: true},
		"Error when passphrases do not match":   {confirm: "some-other-passphrase", wantErr: true},
		"Error when check calls to snapd fails": {checkPassphraseError: true, wantErr: true},
		"Error when response not ok":            {passphraseNotOK: true, wantErr: true},
		"Error when low entropy":                {passphraseLowEntropy: true, wantErr: true},
		"Error when invalid passphrase":         {passphraseInvalid: true, wantErr: true},
		"Error when unsupported":                {passphraseUnsupported: true, wantErr: true},
		"Error when unknown error":              {passphraseUnknownError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				CheckPassphraseError:   tc.checkPassphraseError,
				PassphraseLowEntropy:   tc.passphraseLowEntropy,
				PassphraseInvalid:      tc.passphraseInvalid,
				PassphraseUnsupported:  tc.passphraseUnsupported,
				PassphraseUnknownError: tc.passphraseUnknownError,
				PassphraseNotOK:        tc.passphraseNotOK,
			})

			// Default passphrase if empty
			passphrase := tc.passphrase
			if !tc.wantErr && passphrase == "" {
				passphrase = "my-secure-passphrase"
			}

			// Default confirm to passphrase for success cases
			confirm := tc.confirm
			if !tc.wantErr && confirm == "" {
				confirm = passphrase
			}

			err := tpm.IsValidPassphrase(ctx, mockClient, passphrase, confirm)

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
		})
	}
}

func TestIsValidPIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pin     string
		confirm string

		checkPINError  bool
		pinLowEntropy  bool
		pinInvalid     bool
		pinUnsupported bool
		pinNotOK       bool

		wantErr bool
	}{
		"Success": {},

		"Error when PIN empty":               {wantErr: true},
		"Error when PIN contains non digits": {pin: "12a bc6", wantErr: true},
		"Error when PINs do not match":       {confirm: "654321", wantErr: true},
		"Error when snapd down":              {checkPINError: true, wantErr: true},
		"Error when response not ok":         {pinNotOK: true, wantErr: true},
		"Error when low entropy":             {pinLowEntropy: true, wantErr: true},
		"Error when invalid PIN":             {pinInvalid: true, wantErr: true},
		"Error when unsupported":             {pinUnsupported: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				CheckPINError:  tc.checkPINError,
				PINLowEntropy:  tc.pinLowEntropy,
				PINInvalid:     tc.pinInvalid,
				PINUnsupported: tc.pinUnsupported,
				PINNotOK:       tc.pinNotOK,
			})

			// Default PIN to 123456 if empty
			pin := tc.pin
			if !tc.wantErr && pin == "" {
				pin = "123456"
			}

			// Default confirm to pin for success cases
			confirm := tc.confirm
			if !tc.wantErr && confirm == "" {
				confirm = pin
			}

			err := tpm.IsValidPIN(ctx, mockClient, pin, confirm)

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
		})
	}
}

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
