package main_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/creack/pty"
	"github.com/matryer/is"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	commands := []string{
		"add-passphrase",
		"add-pin",
	}

	tests := map[string]struct {
		nonRootUser bool

		wantErr bool
	}{
		"Success": {},

		"Fail_on_user_privilege": {nonRootUser: true, wantErr: true},
		"Fail_wrong_auth_mode":   {wantErr: true},
		"Fail_on_validating":     {wantErr: true},
		"Fail_on_adding":         {wantErr: true},
	}
	for _, command := range commands {
		for name, tc := range tests {
			t.Run(filepath.Join(command, name), func(t *testing.T) {
				t.Parallel()

				is := is.New(t)

				input := "test"
				if command == "add-pin" {
					input = "12345"
				}

				user := testutils.WithUserAsRoot()
				if tc.nonRootUser {
					user = testutils.WithUserAsNonRoot()
				}

				root, err := filepath.Abs(testutils.TestPath(t))
				is.NoErr(err)

				ptmx, tty, err := pty.Open()
				is.NoErr(err)
				defer ptmx.Close()
				defer tty.Close()

				go func() {
					for range 2 {
						fmt.Fprintln(ptmx, input)
					}
				}()

				//nolint:gosec // The test intentionally executes the binary built in TestMain.
				cmd := exec.Command(cmdPath, command)
				cmd.Env = append(cmd.Env, testutils.WithRootDir(root), user)
				cmd.Stdin = tty

				out, err := cmd.CombinedOutput()
				if testutils.CheckError(is, err, tc.wantErr) {
					return
				}

				golden.CheckOrUpdate(t, out) // TestAdd retruns the correct output
			})
		}
	}
}
