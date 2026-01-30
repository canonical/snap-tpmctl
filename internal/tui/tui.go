// Package tui provides text user interface utilities for interactive terminal operations.
package tui

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/snapcore/snapd/progress"
)

var stdin io.Reader = os.Stdin
var stdout io.Writer = os.Stdout

// these are the bits of the ANSI escapes (beyond \r) that we use
// (names of the terminfo capabilities, see terminfo(5)).
var (
	// clear to end of line.
	clrEOL = "\033[K"
	// make cursor invisible.
	cursorInvisible = "\033[?25l"
	// make cursor visible.
	cursorVisible = "\033[?25h"
)

// ReadUserInput reads a line of input from the user via stdin.
func ReadUserInput() (string, error) {
	reader := bufio.NewReader(stdin)
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	key = strings.TrimSpace(key)

	return key, nil
}

// ClearLine clears the current line in the terminal.
func ClearLine() {
	fmt.Fprint(stdout, "\r", clrEOL)
}

// HideCursor hides the cursor in the terminal.
func HideCursor() {
	fmt.Fprint(stdout, "\r", cursorInvisible, clrEOL)
}

// ShowCursor makes the cursor visible in the terminal.
func ShowCursor() {
	fmt.Fprint(stdout, "\r", cursorVisible, clrEOL)
}

// ReadUserSecret prompts the user for sensitive input and clears the line after reading.
func ReadUserSecret(form string) (string, error) {
	fmt.Fprintf(stdout, "%s", form)
	defer ClearLine()

	input, err := ReadUserInput()
	if err != nil {
		return "", err
	}

	return input, nil
}

// WithSpinner executes an error-only function while displaying a spinner in the terminal.
func WithSpinner(message string, fn func() error) error {
	_, err := WithSpinnerResult(message, func() (struct{}, error) {
		return struct{}{}, fn()
	})
	return err
}

// WithSpinnerResult executes a function while displaying a spinner in the terminal.
func WithSpinnerResult[T any](message string, fn func() (T, error)) (T, error) {
	// Generic result channel
	done := make(chan struct {
		result T
		err    error
	}, 1)

	// Start the func in a goroutine
	go func() {
		result, err := fn()
		done <- struct {
			result T
			err    error
		}{result, err}
	}()

	// Timer to trigger changing the spinner char to produce a loading spinner
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Hide cursor while spinning
	HideCursor()

	var spinner progress.ANSIMeter
	for {
		select {
		case res := <-done:
			spinner.Finished()
			return res.result, res.err
		case <-ticker.C:
			spinner.Spin(message)
		}
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
