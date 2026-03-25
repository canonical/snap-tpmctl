package tpm_test

import (
	"testing"

	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/matryer/is"
)

func TestFdeStatus(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		want    string
		wantErr bool
	}{
		"Returns_FDE_status": {want: "enabled"},

		"Error_when_getting_FDE_status": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			c := snapdtestutils.NewMockSnapdServer(t, ctx)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			got, err := s.FdeStatus(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.Equal(got, tc.want) // TestFDEStatus returns the expected FDE status
		})
	}
}

func TestListVolumeInfo(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Returns_volume_info": {},

		"Error_when_getting_volume_info": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			c := snapdtestutils.NewMockSnapdServer(t, ctx)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			got, err := s.ListVolumeInfo(ctx)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			golden.CheckOrUpdateYAML(t, got) // TestListVolumeInfo returns the expected volume info
		})
	}
}
