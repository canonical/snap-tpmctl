// Package tui provides text user interface utilities for interactive terminal operations.
package tui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/snapcore/snapd/progress"
	"golang.org/x/term"
)

var stdout io.Writer = os.Stdout

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

// ClearPreviousLines clears the previous lines in the terminal.
func ClearPreviousLines(lines int) {
	clr := fmt.Sprint("\r", cursorUp, clrEOL)
	fmt.Fprint(stdout, strings.Repeat(clr, lines))
}

// HideCursor hides the cursor in the terminal.
func HideCursor() {
	fmt.Fprint(stdout, "\r", cursorInvisible, clrEOL)
}

// ShowCursor makes the cursor visible in the terminal.
func ShowCursor() {
	fmt.Fprint(stdout, "\r", cursorVisible, clrEOL)
}

// ReadUserSecret prompts the user for sensitive input with no echo on typing.
func ReadUserSecret(form string) (string, error) {
	fmt.Fprintf(stdout, "%s", form)

	input, err := term.ReadPassword(syscall.Stdin)
	fmt.Fprintln(stdout)

	if err != nil {
		return "", fmt.Errorf("failed to read input: %v", err)
	}

	return string(input), nil
}

// Spin provides a simple interface to start and stop a spinner in the terminal.
func Spin(msg string) (stop func()) {
	var spinner progress.ANSIMeter
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Go(func() {
		// Timer to trigger changing the spinner char to produce a loading spinner
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		// Hide cursor while spinning
		HideCursor()
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
func DisplayTable(w io.Writer, headers []string, rows [][]string, hideHeaders bool) error {
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

	if _, err := io.Copy(w, &buf); err != nil {
		return err
	}

	return nil
}
