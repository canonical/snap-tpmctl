package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func (a App) newCreateKeyCmd() *cli.Command {
	var recoveryKeyName string

	return &cli.Command{
		Name:  "create-recovery-key",
		Usage: "Create a new recovery key",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "key-id",
				UsageText:   "<key-id>",
				Destination: &recoveryKeyName,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Validate the recovery key name
			if err := a.tpm.ValidateRecoveryKeyNameUnique(ctx, recoveryKeyName); err != nil {
				return err
			}

			stop := a.tui.Spin("Generating recovery key...")
			defer stop()

			recoveryKey, err := a.tpm.CreateKey(ctx, recoveryKeyName)
			if err != nil {
				return err
			}

			stop()

			fmt.Printf("Recovery Key: %s\n", recoveryKey)

			// Wait for user to confirm by pressing Enter
			fmt.Print("Save the recovery key somewhere safe. Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			a.tui.ClearPreviousLines(2)

			return nil
		},
	}
}
