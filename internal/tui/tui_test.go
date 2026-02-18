package tui_test

import (
	"errors"
	"testing"
	"time"

	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/matryer/is"
)

func createTestFunc(t *testing.T, wantErr bool) func() error {
	t.Helper()

	if wantErr {
		return func() error {
			return errors.New("operation failed")
		}
	}

	return func() error {
		time.Sleep(250 * time.Millisecond)
		return nil
	}
}

func TestWithSpinner(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Function completes":        {},
		"Error when function fails": {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			fn := createTestFunc(t, tc.wantErr)

			err := tui.WithSpinner("Testing", fn)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
		})
	}
}

func createTestFuncResult(t *testing.T, wantErr bool) func() (string, error) {
	t.Helper()

	if wantErr {
		return func() (string, error) {
			return "", errors.New("operation failed")
		}
	}

	return func() (string, error) {
		time.Sleep(250 * time.Millisecond)
		return "success", nil
	}
}

func TestWithSpinnerResult(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		wantErr bool
	}{
		"Function completes and returns a result": {},
		"Error when function fails":               {wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			fn := createTestFuncResult(t, tc.wantErr)

			val, err := tui.WithSpinnerResult("Testing", fn)
			if tc.wantErr {
				is.Equal("operation failed", err.Error())
				return
			}
			is.NoErr(err)
			is.Equal("success", val)
		})
	}
}
