package tpm_test

import (
	"context"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/matryer/is"
)

//nolint:dupl // CreateKey and RegenerateKey have intentionally similar structure.
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
			is := is.New(t)

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				GenerateKeyError:    tc.generateKeyFails,
				AddRecoveryKeyError: tc.addKeyFails,
			})

			res, err := tpm.CreateKey(ctx, mockClient, tc.recoveryKeyName)

			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			is.Equal("test-key-id-12345", res.KeyID)
			is.Equal("12345-67890-12345-67890-12345-67890-12345-67890", res.RecoveryKey)
		})
	}
}

//nolint:dupl // CreateKey and RegenerateKey have intentionally similar structure.
func TestRegenerateKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKeyName string

		generateKeyFails bool
		replaceKeyFails  bool

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
		"Error when replace key fails": {
			recoveryKeyName: "my-key",
			replaceKeyFails: true,
			wantErr:         true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				GenerateKeyError:        tc.generateKeyFails,
				ReplaceRecoveryKeyError: tc.replaceKeyFails,
			})

			res, err := tpm.RegenerateKey(ctx, mockClient, tc.recoveryKeyName)

			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			is.Equal("test-key-id-12345", res.KeyID)
			is.Equal("12345-67890-12345-67890-12345-67890-12345-67890", res.RecoveryKey)
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
			is := is.New(t)

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(testutils.MockConfig{
				CheckRecoveryKeyError: tc.checkError,
				RecoveryKeyValid:      tc.recoveryKeyValid,
			})

			valid, err := tpm.CheckKey(ctx, mockClient, "test-key")

			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			is.Equal(tc.wantValid, valid)
		})
	}
}
