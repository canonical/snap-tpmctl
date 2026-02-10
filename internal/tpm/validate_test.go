package tpm_test

import (
	"context"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/nalgeon/be"
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

		wantErr bool
	}{
		"Success": {},

		"Error when passphrase empty":           {wantErr: true},
		"Error when passphrases do not match":   {confirm: "some-other-passphrase", wantErr: true},
		"Error when check calls to snapd fails": {checkPassphraseError: true, wantErr: true},
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

		wantErr bool
	}{
		"Success": {},

		"Error when PIN empty":               {wantErr: true},
		"Error when PIN contains non digits": {pin: "12a bc6", wantErr: true},
		"Error when PINs do not match":       {confirm: "654321", wantErr: true},
		"Error when snapd down":              {checkPINError: true, wantErr: true},
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

func TestIsValidRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		key string

		wantErr   bool
		wantInErr string
	}{
		"valid recovery key": {
			key:     "12345-67890-12345-67890-12345-67890-12345-67890",
			wantErr: false,
		},
		"empty key": {
			key:       "",
			wantErr:   true,
			wantInErr: "recovery key cannot be empty",
		},
		"key with letters": {
			key:       "12345-67890-abcde-67890-12345-67890-12345-67890",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters and hyphens",
		},
		"key too short": {
			key:       "12345-67890-12345",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters and hyphens",
		},
		"key too long": {
			key:       "12345-67890-12345-67890-12345-67890-12345-67890-12345",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters and hyphens",
		},
		"key with wrong separator": {
			key:       "12345_67890_12345_67890_12345_67890_12345_67890",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters and hyphens",
		},
		"key with missing separator": {
			key:       "123456789012345678901234567890123456789012345",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters and hyphens",
		},
		"key with four digits": {
			key:       "1234-67890-12345-67890-12345-67890-12345-67890",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters and hyphens",
		},
		"key with six digits": {
			key:       "123456-67890-12345-67890-12345-67890-12345-67890",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters and hyphens",
		},
		"key with spaces": {
			key:       "12345 67890 12345 67890 12345 67890 12345 67890",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters and hyphens",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tpm.ValidateRecoveryKey(tc.key)

			if tc.wantErr {
				be.Err(t, err, tc.wantInErr)
				return
			}
			be.Err(t, err, nil)
		})
	}
}

func TestValidateRecoveryKeyNameUnique(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKeyName string
		enumerateFails  bool
		wantErr         bool
	}{
		"Success": {
			recoveryKeyName: "my-key",
		},
		"Error when name empty": {
			recoveryKeyName: "",
			wantErr:         true,
		},
		"Error when name starts with snap": {
			recoveryKeyName: "snap-key",
			wantErr:         true,
		},
		"Error when name starts with default": {
			recoveryKeyName: "default-key",
			wantErr:         true,
		},
		"Error when name matches existing recovery Key": {
			recoveryKeyName: "additional-recovery",
			wantErr:         true,
		},
		"Error when enumerate fails": {
			recoveryKeyName: "my-key",
			enumerateFails:  true,
			wantErr:         true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				EnumerateError: tc.enumerateFails,
			})

			err := tpm.ValidateRecoveryKeyNameUnique(ctx, mockClient, tc.recoveryKeyName)

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
		})
	}
}

func TestValidateAuthMode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expectedAuthMode snapd.AuthMode
		mockAuthMode     string
		enumerateFails   bool
		wantErr          bool
	}{
		"Validates passphrase authentication in use": {
			expectedAuthMode: snapd.AuthModePassphrase,
		},
		"Validates PIN authentication in use": {
			expectedAuthMode: snapd.AuthModePin,
		},
		"Validates no authentication in use": {
			expectedAuthMode: snapd.AuthModeNone,
		},
		"Error when enumerate fails": {
			expectedAuthMode: snapd.AuthModePassphrase,
			enumerateFails:   true,
			wantErr:          true,
		},
		"Error when auth mode mismatch": {
			expectedAuthMode: snapd.AuthModePin,
			mockAuthMode:     "passphrase",
			wantErr:          true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// Default mock auth mode to expected if not specified
			mockAuthMode := tc.mockAuthMode
			if mockAuthMode == "" {
				mockAuthMode = string(tc.expectedAuthMode)
			}

			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				EnumerateError: tc.enumerateFails,
				AuthMode:       mockAuthMode,
			})

			err := tpm.ValidateAuthMode(ctx, mockClient, tc.expectedAuthMode)

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
		})
	}
}

func TestValidateDevicePath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path string

		wantErr   bool
		wantInErr string
	}{
		"valid device path": {
			path:    "/dev/null",
			wantErr: false,
		},
		"empty path": {
			path:      "",
			wantErr:   true,
			wantInErr: "device path cannot be empty",
		},
		"non-existent device": {
			path:      "/dev/nonexistent-device-12345",
			wantErr:   true,
			wantInErr: "does not exist",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tpm.ValidateDevicePath(tc.path)

			if tc.wantErr {
				be.Err(t, err, tc.wantInErr)
				return
			}
			be.Err(t, err, nil)
		})
	}
}

func TestValidateDirectoryPath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		dest string

		wantErr   bool
		wantInErr string
	}{
		"valid directory path": {
			dest:    "my-volume",
			wantErr: false,
		},
		"valid directory path with numbers": {
			dest:    "volume-123",
			wantErr: false,
		},
		"valid directory absolute path": {
			dest:    "/mnt/volume",
			wantErr: false,
		},
		"valid directory relative path": {
			dest:    "./mnt/volume",
			wantErr: false,
		},
		"valid directory path with subdirectory": {
			dest:    "folder/subfolder/volume",
			wantErr: false,
		},
		"empty directory path": {
			dest:      "",
			wantErr:   true,
			wantInErr: "directory path cannot be empty",
		},
		"directory path escaping with parent directory": {
			dest:      "../../mnt/vol",
			wantErr:   true,
			wantInErr: "directory path must be a valid absolute or relative path",
		},
		"directory path with excessive parent references": {
			dest:      "../../../../../../../root",
			wantErr:   true,
			wantInErr: "directory path must be a valid absolute or relative path",
		},
		"directory path escaping from subdirectory": {
			dest:      "folder/../../vol",
			wantErr:   true,
			wantInErr: "directory path must be a valid absolute or relative path",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tpm.ValidateDirectoryPath(tc.dest)

			if tc.wantErr {
				be.Err(t, err, tc.wantInErr)
				return
			}
			be.Err(t, err, nil)
		})
	}
}
