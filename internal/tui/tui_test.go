package tui_test

import (
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

func TestSpin(t *testing.T) {
	// Capture spinner output to a buffer, with a mutex to avoid race conditions:
	// https://github.com/golang/go/issues/74352
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
