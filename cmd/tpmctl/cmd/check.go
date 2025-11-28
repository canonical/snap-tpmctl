// Package cmd implements the cli for exposing the cli commands snap-tpmctl supports
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
)

func newCheckCmd() *cli.Command {
	return &cli.Command{
		Name:    "check-recovery-key",
		Usage:   "Check recovery key",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			key, err := readUserInput()
			if err != nil {
				return err
			}

			if err := IsValidRecoveryKey(key); err != nil {
				return err
			}

			return check(ctx, key)
		},
	}
}

func check(ctx context.Context, key string) error {
	c := snapd.NewClient()
	defer c.Close()

	if err := c.LoadAuthFromHome(); err != nil {
		return fmt.Errorf("failed to load auth: %w", err)
	}

	res, err := c.CheckRecoveryKey(ctx, key, nil)
	if err != nil {
		return fmt.Errorf("failed to check recovery key: %w", err)
	}

	msg := "Recovery key does not work"
	if res.IsOK() {
		msg = "Recovery key works"
	}

	fmt.Println(msg)

	return nil
}

func readUserInput() (string, error) {
	fmt.Print("Enter recovery key: ")

	reader := bufio.NewReader(os.Stdin)
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read recovery key: %w", err)
	}
	key = strings.TrimSpace(key)

	// Clear the terminal line to hide the recovery key from stdout
	fmt.Print("\033[1A\033[2K")

	return key, nil
}

// IsValidRecoveryKey checks to see if a recovery key matches expected formatting.
func IsValidRecoveryKey(key string) error {
	if key == "" {
		return fmt.Errorf("recovery key cannot be empty")
	}

	matched, err := regexp.MatchString(`^([0-9]{5}-){7}[0-9]{5}$`, key)
	if err != nil {
		return fmt.Errorf("regex validation error: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid recovery key format: must contain only alphanumeric characters, hyphens, or underscores")
	}

	return nil
}
