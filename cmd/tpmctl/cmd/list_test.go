package cmd_test

import (
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	cmdtestutils "github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd/testutils"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/matryer/is"
)

func TestListAll(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		hideHeaders   bool
		tuiWriteError bool

		wantErr bool
	}{
		"Success_on_getting_keyslots":                 {},
		"Success_on_getting_keyslots_without_headers": {hideHeaders: true},

		"Error_on_getting_keyslots":    {wantErr: true},
		"Error_on_displaying_keyslots": {tuiWriteError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)
			ctx, logs := testutils.TestLoggerWithBuffer(t)

			command := "list-all"

			args := []string{command}
			if tc.hideHeaders {
				args = append(args, "--no-headers")
			}

			var out strings.Builder
			w := testWriter{io.Writer(&out), tc.tuiWriteError}
			tui := tui.New(nil, w)

			c := snapdtestutils.NewMockSnapdServer(t, ctx)
			s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))
			app := cmd.New(
				cmdtestutils.WithSnapTPM(s),
				cmdtestutils.WithArgs(args...),
				cmdtestutils.WithTui(tui),
			)

			err := app.Run(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default

			golden.CheckOrUpdate(t, out.String())
		})
	}
}

func TestListFiltered(t *testing.T) {
	t.Parallel()

	commands := []string{
		"list-passphrases",
		"list-recovery-keys",
		"list-pins",
	}

	tests := map[string]struct {
		wantErr bool
	}{
		"Success_on_getting_keyslots": {},

		"Error_on_getting_keyslots": {wantErr: true},
	}

	for _, command := range commands {
		for name, tc := range tests {
			t.Run(filepath.Join(command, name), func(t *testing.T) {
				t.Parallel()

				is := is.New(t)
				ctx, logs := testutils.TestLoggerWithBuffer(t)

				var out strings.Builder
				w := io.Writer(&out)
				tui := tui.New(nil, w)

				c := snapdtestutils.NewMockSnapdServer(t, ctx)
				s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))
				app := cmd.New(
					cmdtestutils.WithSnapTPM(s),
					cmdtestutils.WithArgs(command),
					cmdtestutils.WithTui(tui),
				)

				err := app.Run(ctx)
				if testutils.CheckError(is, err, tc.wantErr) {
					return
				}

				is.True(logs.Len() == 0) // No logs printed by default

				golden.CheckOrUpdate(t, out.String())
			})
		}
	}
}

type testWriter struct {
	w io.Writer

	wantErr bool
}

func (t testWriter) Write(w []byte) (int, error) {
	if t.wantErr {
		return 0, errors.New("write failed")
	}

	return t.w.Write(w)
}
