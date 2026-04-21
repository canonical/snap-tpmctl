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

		golden.CheckOrUpdate(t, globalBuf.String())
	})
}

func TestReadSecret(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input      string
		wantSecret string

		wantReadErr bool
		wantErr     bool
	}{
		"Success":           {},
		"Success_backspace": {input: "test\bx\n", wantSecret: "tesx"},
		"Success_ctrl_c":    {input: "\x03", wantSecret: ""},

		"Fail_reading_input": {wantReadErr: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			validSecret := "mysecret"
			if tc.input == "" {
				tc.input = validSecret + "\n"
				tc.wantSecret = validSecret
			}

			ptmx, tty, err := pty.Open()
			is.NoErr(err)
			defer ptmx.Close()
			defer tty.Close()

			if tc.wantReadErr {
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
			is.Equal(secret, tc.wantSecret)

			golden.CheckOrUpdate(t, out.String())
		})
	}
}

func TestReadRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		wantKey string

		wantReadErr bool
		wantErr     bool
	}{
		"Success":                   {},
		"Success_with_typed_hyphen": {input: "12345-12345\n", wantKey: "1234512345"},
		"Success_backspace":         {input: "1234\bx2345\n", wantKey: "123x2345"},

		"Fail_reading_input": {wantReadErr: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			// 40-digit key (8 groups of 5 digits)
			validKey := "1234512345123451234512345123451234512345"
			if tc.input == "" {
				tc.input = validKey + "\n"
				tc.wantKey = validKey
			}

			ptmx, tty, err := pty.Open()
			is.NoErr(err)
			defer ptmx.Close()
			defer tty.Close()

			if tc.wantReadErr {
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
			is.Equal(key, tc.wantKey)

			golden.CheckOrUpdate(t, out.String())
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
