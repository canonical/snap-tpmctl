package main_test

import (
	"os"
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
			_, logs := testutils.TestLoggerWithBuffer(t)

			root, err := filepath.Abs(testutils.TestPath(t))
			is.NoErr(err)

			command := "status"

			cmd := exec.Command("go", "run", "-tags=integrationtests", "./cmd/tpmctl", command)
			cmd.Dir = testutils.TestProjectRootPath(is)
			cmd.Env = append(os.Environ(),
				"SNAP_TPMCTL_INTEGRATION_TESTS_ADMIN_EUID=0",
				"SNAP_TPMCTL_INTEGRATION_TESTS_ROOT_DIR="+root,
			)

			out, err := cmd.CombinedOutput()
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default

			golden.CheckOrUpdate(t, out)
		})
	}
}
