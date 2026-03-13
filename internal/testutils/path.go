package testutils

import (
	"path/filepath"
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
