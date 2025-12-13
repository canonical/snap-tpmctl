package tui_test

import (
	"errors"
	"testing"
	"time"

	"github.com/nalgeon/be"
	"snap-tpmctl/internal/tui"
)

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

			var fn func() error
			if tc.wantErr {
				fn = func() error {
					return errors.New("operation failed")
				}
			} else {
				fn = func() error {
					time.Sleep(250 * time.Millisecond)
					return nil
				}
			}

			err := tui.WithSpinner("Testing", fn)
			if tc.wantErr {
				be.Err(t, err)
			} else {
				be.Err(t, err, nil)
			}
		})
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

			var fn func() (string, error)
			if tc.wantErr {
				fn = func() (string, error) {
					return "", errors.New("operation failed")
				}
			} else {
				fn = func() (string, error) {
					time.Sleep(250 * time.Millisecond)
					return "success", nil
				}
			}

			val, err := tui.WithSpinnerResult("Testing", fn)
			if tc.wantErr {
				be.Err(t, err)
			} else {
				be.Err(t, err, nil)
				be.Equal(t, "success", val)
			}
		})
	}
}
