package tui_test

import (
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
	_ "unsafe" // Required for go:linkname directives

	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/matryer/is"
)

//go:linkname spinnerStdout github.com/snapcore/snapd/progress.stdout
var spinnerStdout io.Writer

type syncBuffer struct {
	mu  sync.Mutex
	buf strings.Builder
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.String()
}

func (b *syncBuffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf.Reset()
}

func TestSpin(t *testing.T) {
	// Capture spinner output to a buffer.
	var instantBuf, globalBuf syncBuffer

	w := io.MultiWriter(&instantBuf, &globalBuf)

	spinnerStdout = w
	defer func() { spinnerStdout = os.Stdout }()

	synctest.Test(t, func(t *testing.T) {
		is := is.New(t)

		msg := "Some message..."

		stop := tui.Spin(msg)
		defer stop()
		synctest.Wait()

		is.Equal(instantBuf.String(), "")

		time.Sleep(time.Nanosecond)

		for _, sep := range []string{"/", "-", "\\", "|"} {
			time.Sleep(100 * time.Millisecond)
			synctest.Wait()

			is.True(strings.Contains(instantBuf.String(), msg)) // message is present
			is.True(strings.Contains(instantBuf.String(), sep)) // separator progressed

			instantBuf.Reset()
		}

		stop()
		synctest.Wait()

		golden.CheckOrUpdate(t, globalBuf.String())
	})
}

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
