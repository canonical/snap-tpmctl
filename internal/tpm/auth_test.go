package tpm_test

import (
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
	"github.com/canonical/snap-tpmctl/internal/snapd"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/matryer/is"
)

func TestAddPassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		passphrase string

		wantErr bool
	}{
		"Success_adding_passphrase": {passphrase: "test"},

		"Fail_adding_passphrase": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePlatformKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.AddPassphrase(ctx, tc.passphrase)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.True(tpmtestutils.HasBodyContent(is, *c.Requests, tc.passphrase))
		})
	}

}

func TestReplacePassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		old string
		new string

		wantErr bool
	}{
		"Success_replacing_passphrase": {old: "old", new: "new"},

		"Fail_replacing_passphrase": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePassphrase")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.ReplacePassphrase(ctx, tc.old, tc.new)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.True(tpmtestutils.HasBodyContent(is, *c.Requests, tc.old, tc.new))
		})
	}
}

func TestRemovePassphrase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Success_removing_passphrase": {},

		"Fail_removing_passphrase": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePlatformKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.RemovePassphrase(ctx)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.True(tpmtestutils.HasBodyContent(is, *c.Requests, string(snapd.AuthModeNone)))
		})
	}

}

func TestAddPIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pin string

		wantErr bool
	}{
		"Success_adding_pin": {pin: "123456"},

		"Fail_adding_pin": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePlatformKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.AddPIN(ctx, tc.pin)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.True(tpmtestutils.HasBodyContent(is, *c.Requests, tc.pin))
		})
	}
}

func TestReplacePIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		old string
		new string

		wantErr bool
	}{
		"Success_replacing_pin": {old: "123456", new: "654321"},

		"Fail_replacing_pin": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePIN")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.ReplacePIN(ctx, tc.old, tc.new)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.True(tpmtestutils.HasBodyContent(is, *c.Requests, tc.old, tc.new))
		})
	}
}

func TestRemovePIN(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Success_removing_pin": {},

		"Fail_removing_pin": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := log.WithLoggerInContext(t.Context(), t.Output())
			path := tpmtestutils.GetTestPath(t, tc.wantErr, "ReplacePlatformKey")

			c := snapdtestutils.NewMockSnapdServer(t, ctx, path)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

			err := s.RemovePIN(ctx)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.True(tpmtestutils.HasBodyContent(is, *c.Requests, string(snapd.AuthModeNone)))
		})
	}
}
