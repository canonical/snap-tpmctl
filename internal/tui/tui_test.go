package tui_test

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
	_ "unsafe" // Required for go:linkname directives

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/creack/pty"
	"github.com/matryer/is"
	"golang.org/x/term"
)

//go:linkname spinnerStdout github.com/snapcore/snapd/progress.stdout
var spinnerStdout io.Writer

func TestSpin(t *testing.T) {
	// Capture spinner output to a buffer, with a mutex to avoid race conditions:
	// https://github.com/golang/go/issues/74352
	var instantBuf, globalBuf syncBuffer

	w := io.MultiWriter(&instantBuf, &globalBuf)

	spinnerStdout = w
	defer func() { spinnerStdout = os.Stdout }()

	escapes := getEscapes(t)

	tui := tui.New(nil, w)
	synctest.Test(t, func(t *testing.T) {
		is := is.New(t)

		msg := "Some message..."

		stop := tui.Spin(msg)
		defer stop()
		synctest.Wait()

		// the buffer currently contains only the escapes for hiding the cursor
		is.Equal(instantBuf.String(), escapes)

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

		golden.CheckOrUpdate(t, globalBuf.String()) // TestSpin returns the expected spinner output
	})
}

func TestReadSecret(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input        string
		ttyReadError bool

		wantErr bool
	}{
		"Success":                    {},
		"Success_backspace":          {input: "test\bx\n"},
		"Success_ctrl_c":             {input: "\x03"},
		"Success_ignoring_backspace": {input: "\b\b\b\n"},

		"Error_reading_input": {ttyReadError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			validSecret := "mysecret"
			if tc.input == "" {
				tc.input = validSecret + "\n"
			}

			ptmx, tty, err := pty.Open()
			is.NoErr(err) // Setup: could not create fake terminal
			defer ptmx.Close()
			defer tty.Close()

			// Put the TTY in raw mode before writing, so control characters
			// (e.g. \x03) are not consumed by the before readMaskedInput calls MakeRaw.
			//nolint:gosec // This is used on in tests
			_, err = term.MakeRaw(int(tty.Fd()))
			is.NoErr(err) // Setup: could not put tty in raw mode

			if tc.ttyReadError {
				tty = nil
			}

			var out strings.Builder
			tt := tui.New(tty, &out)

			done := make(chan struct{})
			go func() {
				defer close(done)
				fmt.Fprint(ptmx, tc.input)
			}()

			secret, err := tt.ReadUserSecret("Enter passphrase: ")
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}
			is.NoErr(err)

			got := struct {
				Out    string
				Secret string
			}{
				Out:    out.String(),
				Secret: secret,
			}

			golden.CheckOrUpdate(t, got) // TestReadSecret returns the expected output and secret
		})
	}
}

func TestReadRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input        string
		ttyReadError bool

		wantErr bool
	}{
		"Success":                       {},
		"Success_with_typed_hyphen":     {input: "12345-12345\n"},
		"Success_backspace":             {input: "1234\bx2345\n"},
		"Success_removing_separator":    {input: "123451\b12345\n"},
		"Success_ignoring_larger_input": {input: strings.Repeat("12345", 10) + "\n"},

		"Error_reading_input": {ttyReadError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			// 40-digit key (8 groups of 5 digits)
			validKey := "1234512345123451234512345123451234512345"
			if tc.input == "" {
				tc.input = validKey + "\n"
			}

			ptmx, tty, err := pty.Open()
			is.NoErr(err) // Setup: could not create fake terminal
			defer ptmx.Close()
			defer tty.Close()

			if tc.ttyReadError {
				tty = nil
			}

			var out strings.Builder
			tt := tui.New(tty, &out)

			done := make(chan struct{})
			go func() {
				defer close(done)
				fmt.Fprint(ptmx, tc.input)
			}()

			key, err := tt.ReadRecoveryKey()

			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}
			is.NoErr(err)
			is.True(len(key) <= tui.MaxInputLen)

			got := struct {
				Out string
				Key string
			}{
				Out: out.String(),
				Key: key,
			}

			golden.CheckOrUpdate(t, got) // TestReadRecoveryKey returns the expected output
		})
	}
}

func getEscapes(t *testing.T) string {
	t.Helper()

	escapes := strings.Builder{}
	tt := tui.New(nil, &escapes)
	tt.HideCursor()

	return escapes.String()
}

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
