package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/tpm"
)

func newMountVolumeCmd() *cli.Command {
	var devicePath, volumeName string

	return &cli.Command{
		Name:    "mount-volume",
		Usage:   "Unlock and mount a LUKS encrypted volume",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "device-path",
				UsageText:   "<device-path>",
				Destination: &devicePath,
			},
			&cli.StringArg{
				Name:        "volume-name",
				UsageText:   "<volume-name>",
				Destination: &volumeName,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// TODO: add validator for devicePath and volumeName

			if err := tpm.MountVolume(volumeName, devicePath); err != nil {
				return err
			}

			return nil
		},
	}
}

func newUnmountVolumeCmd() *cli.Command {
	var volumeName string

	return &cli.Command{
		Name:    "unmount-volume",
		Usage:   "Unmount and lock a LUKS encrypted volume",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "volume-name",
				UsageText:   "<volume-name>",
				Destination: &volumeName,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// TODO: add validator for volumeName

			if err := tpm.UnmountVolume(volumeName); err != nil {
				return err
			}

			return nil
		},
	}
}

func newGetLuksPassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:    "get-luks-passphrase",
		Usage:   "Get LUKS passphrase from recovery key",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
}
