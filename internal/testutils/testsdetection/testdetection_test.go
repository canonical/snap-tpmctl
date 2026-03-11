package testsdetection_test

import (
	"os/exec"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
	"github.com/matryer/is"
)

func TestMustBeTestingInTests(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	defer func() {
		r := recover()
		is.Equal(r, nil) // MustBeTesting should not panic as we are running in tests
	}()

	testsdetection.MustBeTesting()
}

func TestMustBeTestingForIntegrationTests(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		integrationtestsTag bool

		wantPanic bool
	}{
		"Pass_when_called_in_an_integration_tests_build": {integrationtestsTag: true, wantPanic: false},

		"Panics_when_called_in_non_tests_and_no_integration_tests": {integrationtestsTag: false, wantPanic: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			args := []string{"run"}
			if tc.integrationtestsTag {
				args = append(args, "-tags=integrationtests")
			}
			// TODO
			/*if testutils.CoverDirForTests() != "" {
				args = append(args, "-cover")
			}
			if testutils.IsRace() {
				args = append(args, "-race")
			}*/
			args = append(args, "testdata/binary.go")

			// Execute our subprocess
			cmd := exec.Command("go", args...)
			//cmd.Env = testutils.AppendCovEnv(os.Environ())
			_, err := cmd.CombinedOutput()

			if tc.wantPanic {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
		})
	}
}
