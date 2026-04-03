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

func newRegenerateKeyCmd() *cli.Command {
	var recoveryKeyName string

	return &cli.Command{
		Name:    "regenerate-recovery-key",
		Usage:   "Regenerate an existing recovery key",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "key-id",
				UsageText:   "<key-id>",
				Destination: &recoveryKeyName,
			},
		},
		ShellComplete: func(ctx context.Context, cmd *cli.Command) {
			if alreadyCompleted(cmd) {
				return
			}

			c := snapd.New()

			result, err := c.ListVolumeInfo(ctx)
			if err != nil {
				return
			}

			data := parseKeySlots(result, snapd.IsRecoveryKey)
			for _, name := range data {
				fmt.Fprintf(cmd.Root().Writer, "%s\n", name)
			}
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			s := tpm.New()

			// Validate the recovery key name
			if err := tpm.ValidateRecoveryKeyName(ctx, recoveryKeyName); err != nil {
				return err
			}

			stop := tui.Spin("Regenerating recovery key...")
			defer stop()

			recoveryKey, err := s.RegenerateKey(ctx, recoveryKeyName)
			if err != nil {
				return err
			}
			stop()

			fmt.Printf("Recovery Key: %s\n", recoveryKey)

			// Wait for user to confirm by pressing Enter
			fmt.Print("Save the recovery key somewhere safe. Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			tui.ClearPreviousLines(2)

			return nil
		},
	}
}
