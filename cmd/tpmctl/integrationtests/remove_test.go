package main_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/matryer/is"
)

func TestRemove(t *testing.T) {
	t.Parallel()

	commands := []string{
		"remove-passphrase",
		"remove-pin",
	}

	tests := map[string]struct {
		nonRootUser bool

		wantErr bool
	}{
		"Success_on_removing": {},

		"Fail_on_user_privilege": {nonRootUser: true, wantErr: true},
		"Fail_on_removing":       {wantErr: true},
	}
	for _, command := range commands {
		for name, tc := range tests {
			t.Run(filepath.Join(command, name), func(t *testing.T) {
				t.Parallel()

				is := is.New(t)

				user := testutils.WithUserAsRoot()
				if tc.nonRootUser {
					user = testutils.WithUserAsNonRoot()
				}

				root, err := filepath.Abs(testutils.TestPath(t))
				is.NoErr(err) // Setup: could not find test path

				//nolint:gosec // The test intentionally executes the binary built in TestMain.
				cmd := exec.Command(cmdPath, command)
				cmd.Env = append(cmd.Env, testutils.WithRootDir(root), user)

				out, err := cmd.CombinedOutput()
				if testutils.CheckError(is, err, tc.wantErr) {
					return
				}

				golden.CheckOrUpdate(t, out) // TestRemove returns the correct output
			})
		}
	}
}
