package main_test

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/matryer/is"
)

var cmdPath string

func TestRun(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	command := "help"

	root, err := filepath.Abs(testutils.TestPath(t))
	is.NoErr(err)

	cmd := exec.Command(cmdPath, command)
	cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsNonRoot())

	err = cmd.Run()
	is.NoErr(err) // Expected no error from this command
}

func TestVersion(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	command := "version"

	root, err := filepath.Abs(testutils.TestPath(t))
	is.NoErr(err)

	cmd := exec.Command(cmdPath, command)
	cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsNonRoot())

	err = cmd.Run()
	is.NoErr(err) // Expected no error from this command
}

func TestMain(m *testing.M) {
	var cleanup func()
	var err error

	cmdPath, cleanup, err = testutils.BuildSnapTpmCtl()
	if err != nil {
		log.Printf("Setup: failed to build snap-tpmctl: %v", err)
		os.Exit(1)
	}
	defer cleanup()

	m.Run()
}
