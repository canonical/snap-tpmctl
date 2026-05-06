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

func TestReplace(t *testing.T) {
	t.Parallel()

	commands := []string{
		"replace-passphrase",
		"replace-pin",
	}

	tests := map[string]struct {
		wantErr bool
	}{
		"Success_on_replacing": {},

		"Fail_on_validating": {wantErr: true},
		"Fail_on_replacing":  {wantErr: true},
	}

	for _, command := range commands {
		for name, tc := range tests {
			t.Run(filepath.Join(command, name), func(t *testing.T) {
				t.Parallel()

				is := is.New(t)

				input := "test"
				if command == "replace-pin" {
					input = "12345"
				}

				root, err := filepath.Abs(testutils.TestPath(t))
				is.NoErr(err) // Setup: could not find test root

				ptmx, tty, err := pty.Open()
				is.NoErr(err) // Setup: could not create fake terminal
				defer ptmx.Close()
				defer tty.Close()

				go func() {
					for range 3 {
						fmt.Fprintln(ptmx, input)
					}
				}()

				//nolint:gosec // The test intentionally executes the binary built in TestMain.
				cmd := exec.Command(cmdPath, command)
				cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsNonRoot())
				cmd.Stdin = tty

				out, err := cmd.CombinedOutput()
				if testutils.CheckError(is, err, tc.wantErr) {
					return
				}

				golden.CheckOrUpdate(t, out) // TestReplace returns the correct output
			})
		}
	}
}
