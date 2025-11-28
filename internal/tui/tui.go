// Package tui provides text user interface utilities for interactive terminal operations.
package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
