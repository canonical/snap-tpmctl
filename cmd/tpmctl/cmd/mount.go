package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/tpm"
	"snap-tpmctl/internal/tui"
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
			if err := tpm.ValidateDevicePath(devicePath); err != nil {
				return err
			}

			if err := tpm.ValidateVolumeName(volumeName); err != nil {
				return err
			}

			// TODO: enable this when snap can mount luks_crypt volumes
			// if err := tpm.MountVolume(volumeName, devicePath); err != nil {
			// 	return err
			// }

			recoveryKey, err := tui.ReadUserSecret("Enter recovery key: ")
			if err != nil {
				return err
			}

			key, err := tpm.GetLuksKey(recoveryKey)
			if err != nil {
				return err
			}

			keyString := strings.ReplaceAll(fmt.Sprintf("%q", key), "\"", "'")

			fmt.Printf("To mount %[1]s as %[2]s, use this command:\n\nprintf %[3]s | systemd-cryptsetup attach %[2]s %[1]s /dev/stdin luks,keyslot=-1,tries=1\n", devicePath, volumeName, keyString)

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
			if err := tpm.ValidateVolumeName(volumeName); err != nil {
				return err
			}

			if err := tpm.UnmountVolume(volumeName); err != nil {
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
