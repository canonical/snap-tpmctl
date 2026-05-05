package testutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

// TestFamilyPath returns the path of the dir for storing fixtures and other files related to the test.
func TestFamilyPath(t *testing.T) string {
	t.Helper()

	// Ensures that only the name of the top level test is used
	topLevelTest, _, _ := strings.Cut(t.Name(), "/")

	return filepath.Join("testdata", topLevelTest)
}

// TestPath returns the path based on the current test name.
func TestPath(t *testing.T) string {
	t.Helper()

	return filepath.Join("testdata", t.Name())
}

// TestProjectRootPath returns the repository root path by walking up from the current working directory until go.mod is found.
func TestProjectRootPath(is *is.I) string {
	is.Helper()

	wd, err := os.Getwd()
	is.NoErr(err)

	dir := filepath.Clean(wd)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		is.True(parent != dir) // Setup: root directory not found

		dir = parent
	}
}
