package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/urfave/cli/v3"
)

func newCreateKeyCmd() *cli.Command {
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
			c := snapd.NewClient()

			// Validate the recovery key name
			if err := tpm.ValidateRecoveryKeyNameUnique(ctx, c, recoveryKeyName); err != nil {
				return err
			}

			result, err := tui.WithSpinnerResult("Generating recovery key...", func() (tpm.CreateKeyResult, error) {
				return tpm.CreateKey(ctx, c, recoveryKeyName)
			})
			if err != nil {
				return err
			}

			fmt.Printf("Recovery Key: %s\n", result.RecoveryKey)

			// Wait for user to confirm by pressing Enter
			fmt.Print("Save the recovery key somewhere safe. Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			tui.ClearPreviousLines(2)

			return nil
		},
	}
}
