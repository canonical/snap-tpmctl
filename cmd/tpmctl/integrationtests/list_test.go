package main_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/matryer/is"
)

func TestListAll(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		hideHeaders bool

		wantErr bool
	}{
		"Success_on_getting_keyslots":                 {},
		"Success_on_getting_keyslots_without_headers": {hideHeaders: true},

		"Error_on_getting_keyslots": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			command := "list-all"

			args := []string{command}
			if tc.hideHeaders {
				args = append(args, "--no-headers")
			}

			root, err := filepath.Abs(testutils.TestPath(t))
			is.NoErr(err) // Setup: could not find test path

			cmd := exec.Command(cmdPath, args...)
			cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsRoot())

			out, err := cmd.CombinedOutput()
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			golden.CheckOrUpdate(t, out)
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

				root, err := filepath.Abs(testutils.TestPath(t))
				is.NoErr(err) // Setup: could not find test path

				//nolint:gosec // The test intentionally executes the binary built in TestMain.
				cmd := exec.Command(cmdPath, command)
				cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsRoot())

				out, err := cmd.CombinedOutput()
				if testutils.CheckError(is, err, tc.wantErr) {
					return
				}

				golden.CheckOrUpdate(t, out)
			})
		}
	}
}
