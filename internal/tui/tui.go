// Package tui provides text user interface utilities for interactive terminal operations.
package tui

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/snapcore/snapd/progress"
	"golang.org/x/term"
)

// these are the bits of the ANSI escapes (beyond \r) that we use
// (names of the terminfo capabilities, see terminfo(5)).
var (
	// clear to end of line.
	clrEOL = "\033[K"
	// move cursor up one line.
	cursorUp = "\033[1A"
	// make cursor invisible.
	cursorInvisible = "\033[?25l"
	// make cursor visible.
	cursorVisible = "\033[?25h"
)

// TerminalReader defines the input stream contract required by Tui.
type TerminalReader interface {
	io.Reader
	Fd() uintptr
}

// Tui wraps reader and writer streams used by terminal UI helpers.
type Tui struct {
	r TerminalReader
	w io.Writer
}

// New returns a Tui configured with the provided reader and writer streams.
func New(r TerminalReader, w io.Writer) Tui {
	return Tui{r, w}
}

// Writer returns the output writer configured for this Tui instance.
func (t Tui) Writer() io.Writer {
	return t.w
}

// Reader returns the input reader configured for this Tui instance.
func (t Tui) Reader() io.Reader {
	return t.r
}

// ClearPreviousLines clears the previous lines in the terminal.
func (t Tui) ClearPreviousLines(lines int) {
	clr := fmt.Sprint(t.w, "\r", cursorUp, clrEOL)
	fmt.Fprint(t.w, strings.Repeat(clr, lines))
}

// HideCursor hides the cursor in the terminal.
func (t Tui) HideCursor() {
	fmt.Fprint(t.w, "\r", cursorInvisible, clrEOL)
}

// ShowCursor makes the cursor visible in the terminal.
func (t Tui) ShowCursor() {
	fmt.Fprint(t.w, "\r", cursorVisible, clrEOL)
}

// ReadUserSecret prompts the user for sensitive input with no echo on typing.
func (t Tui) ReadUserSecret(form string) (string, error) {
	fmt.Fprintf(t.w, "%s", form)

	fd := t.r.Fd()
	const maxInt = int(^uint(0) >> 1)
	if fd > uintptr(maxInt) {
		return "", fmt.Errorf("failed to read input: invalid file descriptor")
	}

	input, err := term.ReadPassword(int(fd))
	if err != nil {
		return "", fmt.Errorf("failed to read input: %v", err)
	}
	fmt.Fprintln(t.w)

	return string(input), nil
}

// Spin provides a simple interface to start and stop a spinner in the terminal.
func (t Tui) Spin(msg string) (stop func()) {
	var spinner progress.ANSIMeter
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Go(func() {
		// Timer to trigger changing the spinner char to produce a loading spinner
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		// Hide cursor while spinning
		t.HideCursor()
		for {
			select {
			case <-done:
				spinner.Finished()
				return
			case <-ticker.C:
				spinner.Spin(msg)
			}
		}
	})

	return func() {
		if done == nil {
			return
		}

		close(done)
		wg.Wait()
		done = nil
	}
}

// DisplayTable writes a formatted table to the given writer with optional headers.
func (t Tui) DisplayTable(headers []string, rows [][]string, hideHeaders bool) error {
	if len(rows) == 0 {
		return nil
	}

	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	if !hideHeaders {
		fmt.Fprintln(tw, strings.Join(headers, "\t"))
	}

	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}

	if err := tw.Flush(); err != nil {
		return err
	}

	if _, err := io.Copy(t.w, &buf); err != nil {
		return err
	}

	return nil
}
