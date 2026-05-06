package main_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/matryer/is"
)

func TestStatus(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Returns_FDE_status": {},

		"Error_when_getting_FDE_status": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			command := "status"

			root, err := filepath.Abs(testutils.TestPath(t))
			is.NoErr(err) // Setup: could not find test path

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
