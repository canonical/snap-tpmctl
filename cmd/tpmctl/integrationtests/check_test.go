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

func TestCheck(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		key string

		wantErr bool
	}{
		"Success_checking_recovery_key":            {},
		"Success_even_with_invalid_recovery_key":   {},
		"Success_even_with_incorrect_recovery_key": {key: "incorrect"},

		"Fail_checking_recovery_key": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			command := "check-recovery-key"
			if tc.key == "" {
				tc.key = "11272-47509-28031-54818-41671-38673-11053-06376"
			}

			root, err := filepath.Abs(testutils.TestPath(t))
			is.NoErr(err) // Setup: could not find test root

			ptmx, tty, err := pty.Open()
			is.NoErr(err) // Setup: could not create fake terminal
			defer ptmx.Close()
			defer tty.Close()

			go func() {
				fmt.Fprintln(ptmx, tc.key)
			}()

			cmd := exec.Command(cmdPath, command)
			cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsNonRoot())
			cmd.Stdin = tty

			out, err := cmd.CombinedOutput()
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			golden.CheckOrUpdate(t, out) // TestCheck returns the expected output
		})
	}
}
