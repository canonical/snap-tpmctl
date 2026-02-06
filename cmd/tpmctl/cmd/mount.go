package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/urfave/cli/v3"
)

func newMountVolumeCmd() *cli.Command {
	var device, dir string

	return &cli.Command{
		Name:    "mount-volume",
		Usage:   "Unlock and mount a LUKS encrypted volume",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "device",
				UsageText:   "<device>",
				Destination: &device,
			},
			&cli.StringArg{
				Name:        "dir",
				UsageText:   "<dir>",
				Destination: &dir,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := tpm.ValidateDevicePath(device); err != nil {
				return err
			}

			if err := tpm.ValidateDiretoryPath(dir); err != nil {
				return err
			}

			if err := tpm.MountVolume(dir, device); err != nil {
				return err
			}

			return nil
		},
	}
}

func newUnmountVolumeCmd() *cli.Command {
	var dir string

	return &cli.Command{
		Name:    "unmount-volume",
		Usage:   "Unmount and lock a LUKS encrypted volume",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "dir",
				UsageText:   "<dir>",
				Destination: &dir,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := tpm.ValidateDiretoryPath(dir); err != nil {
				return err
			}

			if err := tpm.UnmountVolume(dir); err != nil {
				return err
			}

			return nil
		},
	}
}

func newGetLuksKeyFromRecoveryKeyCmd() *cli.Command {
	var outputFile string
	var hex, escaped bool

	return &cli.Command{
		Name:    "get-luks-key",
		Usage:   "Get LUKS key from recovery key",
		Suggest: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Usage:       "Write binary key to file with secure permissions (600)",
				Destination: &outputFile,
			},
			&cli.BoolFlag{
				Name:        "hex",
				Usage:       "Output key in hexadecimal format",
				Destination: &hex,
			},
			&cli.BoolFlag{
				Name:        "escaped",
				Usage:       "Output key in escaped string format",
				Destination: &escaped,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			recoveryKey, err := tui.ReadUserSecret("Enter recovery key: ")
			if err != nil {
				return err
			}

			key, err := tpm.GetLuksKey(recoveryKey)
			if err != nil {
				return err
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, key, 0600); err != nil {
					return fmt.Errorf("failed to write key to file: %w", err)
				}
				fmt.Printf("Binary key written to: %s\n", outputFile)

				return nil
			}

			switch {
			case hex:
				fmt.Printf("%x\n", key)
			case escaped:
				fmt.Printf("%q\n", key)
			default:
				fmt.Printf("LUKS key (hex): %x\n", key)
			}

			return nil
		},
	}
}
