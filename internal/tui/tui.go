// Package tui provides text user interface utilities for interactive terminal operations.
package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ReadUserInput reads a line of input from the user via stdin.
func ReadUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	key = strings.TrimSpace(key)

	return key, nil
}

// ClearLine clears the current line in the terminal.
func ClearLine() {
	fmt.Print("\033[1A\033[2K")
}

// clearCurrentLine clears the current line without moving cursor up.
func clearCurrentLine() {
	fmt.Print("\r\033[K")
}

// ReadUserSecret prompts the user for sensitive input and clears the line after reading.
func ReadUserSecret(form string) (string, error) {
	fmt.Print(form)
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
	spinnerChars := []string{"-", "\\", "|", "/"}
	i := 0

	// Timer to trigger changing the spinner char to produce a loading spinner
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

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

	// Hide cursor while spinning
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	// Spin until we get a result from the function
	fmt.Printf("%s %s", message, spinnerChars[0])
	for {
		select {
		case res := <-done:
			clearCurrentLine()
			return res.result, res.err
		case <-ticker.C:
			clearCurrentLine()
			i++
			fmt.Printf("%s %s", message, spinnerChars[i%len(spinnerChars)])
		}
	}
}
