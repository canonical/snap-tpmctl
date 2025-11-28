package cmd_test

import (
	"bytes"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"snap-tpmctl/cmd/tpmctl/cmd"
)

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
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters, hyphens, or underscores",
		},
		"key too short": {
			key:       "12345-67890-12345",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters, hyphens, or underscores",
		},
		"key too long": {
			key:       "12345-67890-12345-67890-12345-67890-12345-67890-12345",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters, hyphens, or underscores",
		},
		"key with wrong separator": {
			key:       "12345_67890_12345_67890_12345_67890_12345_67890",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters, hyphens, or underscores",
		},
		"key with missing separator": {
			key:       "123456789012345678901234567890123456789012345",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters, hyphens, or underscores",
		},
		"key with four digits": {
			key:       "1234-67890-12345-67890-12345-67890-12345-67890",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters, hyphens, or underscores",
		},
		"key with six digits": {
			key:       "123456-67890-12345-67890-12345-67890-12345-67890",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters, hyphens, or underscores",
		},
		"key with spaces": {
			key:       "12345 67890 12345 67890 12345 67890 12345 67890",
			wantErr:   true,
			wantInErr: "invalid recovery key format: must contain only alphanumeric characters, hyphens, or underscores",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var logs bytes.Buffer

			out := io.MultiWriter(&logs, t.Output())
			h := slog.NewTextHandler(out, nil)
			_ = slog.New(h)

			err := cmd.IsValidRecoveryKey(tc.key)
			if tc.wantErr {
				require.Error(t, err, "Expected an error but got none")
				require.Contains(t, err.Error(), tc.wantInErr, "Error message does not contain expected text")
			} else {
				require.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}
