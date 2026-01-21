package tpm_test

import (
	"context"
	"testing"

	"github.com/nalgeon/be"
	"snap-tpmctl/internal/testutils"
	"snap-tpmctl/internal/tpm"
)

func TestFdeStatus(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg     testutils.MockConfig
		want    string
		wantErr bool
	}{
		"Success default value": {
			want: "active",
		},
		"Success custom value": {
			cfg:  testutils.MockConfig{FdeStatusValue: "inactive"},
			want: "inactive",
		},

		"Error from snapd": {
			cfg:     testutils.MockConfig{FdeStatusError: true},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := testutils.NewMockSnapdClient(tc.cfg)

			res, err := tpm.FdeStatus(ctx, mockClient)

			if tc.wantErr {
				be.Err(t, err)
				return
			}

			be.Err(t, err, nil)
			be.Equal(t, tc.want, res)
		})
	}
}
