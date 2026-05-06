package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	adminEuidEnv = "SNAP_TPMCTL_INTEGRATION_TESTS_ADMIN_EUID"
	rootDirEnv   = "SNAP_TPMCTL_INTEGRATION_TESTS_ROOT_DIR"
)

// GetEuidEnv returns the admin EUID configured for integration tests.
// It panics if the environment variable is not a valid integer.
func GetEuidEnv() int {
	euid, err := strconv.Atoi(os.Getenv(adminEuidEnv))
	if err != nil {
		panic(adminEuidEnv + " must be a valid integer")
	}

	return euid
}

// WithUserAsRoot returns an environment assignment string for the admin EUID.
func WithUserAsRoot() string {
	return fmt.Sprintf("%s=%d", adminEuidEnv, 0)
}

// WithUserAsNonRoot returns an environment assignment string for a non-root EUID.
func WithUserAsNonRoot() string {
	return fmt.Sprintf("%s=%d", adminEuidEnv, 1)
}

// GetRootDirEnv returns the root directory configured for integration tests.
// It panics if the environment variable is empty.
func GetRootDirEnv() string {
	root := os.Getenv(rootDirEnv)
	if root == "" {
		panic(rootDirEnv + " must be set")
	}

	return root
}

// WithRootDir returns an environment assignment string for the integration test root directory.
func WithRootDir(path string) string {
	return fmt.Sprintf("%s=%s", rootDirEnv, path)
}

// BuildSnapTpmCtl builds the executable and returns the binary path.
func BuildSnapTpmCtl() (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "snap-tpmctl-tests")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	cleanup := func() { os.RemoveAll(tempDir) }

	execPath := filepath.Join(tempDir, "snap-tpmctl")
	cmd := exec.Command("go", "build")
	cmd.Dir = ProjectRoot()
	cmd.Args = append(cmd.Args, "-tags=integrationtests")
	cmd.Args = append(cmd.Args, "-o", execPath, "./cmd/tpmctl")

	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to build snap-tpmctl: %v", err)
	}

	return execPath, cleanup, err
}
