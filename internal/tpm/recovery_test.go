package tpm_test

import (
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/matryer/is"
)

func TestCreateKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKeyName string
		recoveryKey     string

		wantGenErr bool
		wantAddErr bool
		wantErr    bool
	}{
		"Success_creating_recovery_key": {recoveryKeyName: "test", recoveryKey: "11272-47509-28031-54818-41671-38673-11053-06376"},

		"Fail_generating_recovery_key": {wantGenErr: true, wantErr: true},
		"Fail_adding_recovery_key":     {wantAddErr: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			pathGen := tpmtestutils.GetTestPath(t, tc.wantGenErr, "GenerateRecoveryKey")
			pathAdd := tpmtestutils.GetTestPath(t, tc.wantAddErr, "AddRecoveryKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, pathGen, pathAdd)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			got, err := s.CreateKey(ctx, tc.recoveryKeyName)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			is.Equal(got, tc.recoveryKey)

			is.True(tpmtestutils.HasBodyContent(is, *c.Requests, tc.recoveryKeyName))
		})
	}
}

func TestRegenerateKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKeyName string
		recoveryKey     string

		wantGenErr bool
		wantAddErr bool
		wantErr    bool
	}{
		"Success_regenerating_recovery_key": {recoveryKeyName: "test", recoveryKey: "11272-47509-28031-54818-41671-38673-11053-06376"},

		"Fail_generating_recovery_key": {wantGenErr: true, wantErr: true},
		"Fail_adding_recovery_key":     {wantAddErr: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			pathGen := tpmtestutils.GetTestPath(t, tc.wantGenErr, "GenerateRecoveryKey")
			pathAdd := tpmtestutils.GetTestPath(t, tc.wantAddErr, "ReplaceRecoveryKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, pathGen, pathAdd)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			got, err := s.RegenerateKey(ctx, tc.recoveryKeyName)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			is.Equal(got, tc.recoveryKey)

			is.True(tpmtestutils.HasBodyContent(is, *c.Requests, tc.recoveryKeyName))
		})
	}
}

func TestCheckKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKey string

		wantErr bool
	}{
		"Success_checking_recovery_key": {recoveryKey: "11272-47509-28031-54818-41671-38673-11053-06376"},

		"Fail_checking_recovery_key": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			path := tpmtestutils.GetTestPath(t, tc.wantErr, "CheckRecoveryKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			got, err := s.CheckKey(ctx, tc.recoveryKey)
			if tc.wantErr {
				is.True(err != nil)
				is.Equal(got, false)
				return
			}
			is.NoErr(err)
			is.Equal(got, true)

			is.True(tpmtestutils.HasBodyContent(is, *c.Requests, tc.recoveryKey))
		})
	}
}

func TestGetLuksKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKey string

		wantErr bool
	}{
		"Success_getting_luks_key": {recoveryKey: "11272-47509-28031-54818-41671-38673-11053-06376"},

		"Fail_getting_luks_key": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())

			got, err := tpm.GetLuksKey(ctx, tc.recoveryKey)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			golden.CheckOrUpdate(t, string(got))
		})
	}
}
