// Package golden provides utilities to compare and update golden files in tests.
package golden

import (
	"go/ast"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/matryer/is"
	"gopkg.in/yaml.v3"
)

var update bool

const (
	// updateGoldenFilesEnv is the environment variable used to indicate go test that
	// the golden files should be overwritten with the current test results.
	updateGoldenFilesEnv = `TESTS_UPDATE_GOLDEN`
)

func init() {
	if os.Getenv(updateGoldenFilesEnv) != "" {
		update = true
	}
}

// CheckOrUpdateYAML compares the provided object with the content of the golden file. If the update environment
// variable is set, the golden file is updated with the provided object serialized as YAML.
func CheckOrUpdateYAML[T any](t *testing.T, got T) {
	t.Helper()

	is := is.New(t)
	goldenFile := goldenPath(t)

	if update {
		data, err := yaml.Marshal(got)
		is.NoErr(err) // Golden: cannot serialize provided object
		updateGoldenFile(t, goldenFile, data)
	}

	t.Logf("Comparing with %q", goldenFile)
	src, err := os.ReadFile(goldenFile)
	is.NoErr(err) // Golden: cannot read golden file
	var want T
	err = yaml.Unmarshal(src, &want)
	is.NoErr(err) // Golden: cannot deserialize golden file content

	diff := cmp.Diff(want, got,
		cmpopts.EquateEmpty(), // Treat empty slices and maps as equal to nil
		cmp.FilterPath(func(p cmp.Path) bool {
			if sf, ok := p.Last().(cmp.StructField); ok {
				return !ast.IsExported(sf.Name())
			}
			return false
		}, cmp.Ignore()), // Ignore all unexported fields.
	)
	if diff != "" {
		t.Logf("Difference between golden file and actual output (-want +got):\n%s", diff)
		t.Fatal()
	}
}

// updateGoldenFile updates the golden file at the specified path with the provided data.
func updateGoldenFile(t *testing.T, path string, data []byte) {
	t.Helper()
	is := is.New(t)

	t.Logf("updating golden file %s", path)
	err := os.MkdirAll(filepath.Dir(path), 0750)
	is.NoErr(err) // Golden: cannot create directory required for golden files
	err = os.WriteFile(path, data, 0600)
	is.NoErr(err) // Golden: cannot update golden file
}

// goldenPath returns the golden goldenPath for the provided test after asserting it’s valid.
func goldenPath(t *testing.T) string {
	t.Helper()
	is := is.New(t)

	// Replace below the regexp with the regexp package in go:
	is.True(regexp.MustCompile(`^[\w\-.\/]+$`).MatchString(t.Name())) // Golden: Invalid golden file name. Only alphanumeric characters, underscores, dashes, and dots are allowed

	return filepath.Join("testdata", "golden", t.Name())
}
