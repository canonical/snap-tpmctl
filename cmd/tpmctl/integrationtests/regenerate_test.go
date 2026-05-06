package main_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/matryer/is"
)

func TestRegenerateKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKeyName string

		wantErr bool
	}{
		"Success_on_regenerating_recovery_key": {},

		"Error_on_empty_name":                {wantErr: true},
		"Error_on_regenerating_recovery_key": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			command := "regenerate-recovery-key"

			if tc.recoveryKeyName == "" {
				tc.recoveryKeyName = "test"
			}
			root, err := filepath.Abs(testutils.TestPath(t))
			is.NoErr(err) // Setup: could not find test path

			//nolint:gosec // The test intentionally executes the binary built in TestMain.
			cmd := exec.Command(cmdPath, command, tc.recoveryKeyName)
			cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsRoot())

			out, err := cmd.CombinedOutput()
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			golden.CheckOrUpdate(t, out) // TestRegenerateKey returns the expected output
		})
	}
}
