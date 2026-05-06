package testutils

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
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

// ProjectRoot returns the absolute path to the project root.
func ProjectRoot() string {
	// p is the path to the current file, in this case -> {PROJECT_ROOT}/internal/testutils/path.go
	_, p, _, _ := runtime.Caller(0)

	// Walk up the tree to get the path of the project root
	l := strings.Split(p, "/")

	// Ignores the last 3 elements -> /internal/testutils/path.go
	l = l[:len(l)-3]

	// strings.Split removes the first "/" that indicated an AbsPath, so we append it back in the final string.
	return "/" + filepath.Join(l...)
}
