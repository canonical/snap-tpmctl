package tpm_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/nalgeon/be"
	"snap-tpmctl/internal/testutils"
	"snap-tpmctl/internal/tpm"
)

func TestCreateKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKeyName string

		generateKeyFails bool
		addKeyFails      bool

		wantErr bool
	}{
		"Success": {
			recoveryKeyName: "my-key",
		},
		"Error when generate key fails": {
			recoveryKeyName:  "my-key",
			generateKeyFails: true,
			wantErr:          true,
		},
		"Error when add key fails": {
			recoveryKeyName: "my-key",
			addKeyFails:     true,
			wantErr:         true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				GenerateKeyError: tc.generateKeyFails,
				AddKeyError:      tc.addKeyFails,
			})

			res, err := tpm.CreateKey(ctx, mockClient, tc.recoveryKeyName)

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
			be.Equal(t, "test-key-id-12345", res.KeyID)
			be.Equal(t, "12345-67890-12345-67890-12345-67890-12345-67890", res.RecoveryKey)
			be.Equal(t, "Done", res.Status)
		})
	}
}

func TestCheckKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		checkError       bool
		recoveryKeyValid bool
		wantValid        bool
		wantErr          bool
	}{
		"Success": {
			recoveryKeyValid: true,
			wantValid:        true,
		},
		"Invalid recovery key": {
			recoveryKeyValid: false,
			wantErr:          true,
		},
		"Check error": {
			checkError: true,
			wantErr:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				CheckRecoveryKeyError: tc.checkError,
				RecoveryKeyValid:      tc.recoveryKeyValid,
			})

			valid, err := tpm.CheckKey(ctx, mockClient, "test-key")

			if tc.wantErr {
				be.Err(t, err)
				return
			}
			be.Err(t, err, nil)
			be.Equal(t, tc.wantValid, valid)
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

			var logs bytes.Buffer

			out := io.MultiWriter(&logs, t.Output())
			h := slog.NewTextHandler(out, nil)
			_ = slog.New(h)

			err := tpm.ValidateRecoveryKey(tc.key)

			if tc.wantErr {
				be.Err(t, err, tc.wantInErr)
				return
			}
			be.Err(t, err, nil)
		})
	}
}
