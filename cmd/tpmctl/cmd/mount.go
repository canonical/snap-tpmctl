package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/snapcore/secboot"
	"github.com/urfave/cli/v3"
)

type authRequestor struct{}

func (r *authRequestor) RequestUserCredential(ctx context.Context, name, path string, authTypes secboot.UserAuthType) (string, error) {
	if authTypes != secboot.UserAuthTypeRecoveryKey {
		return "", fmt.Errorf("authentication type not supported")
	}

	key, err := tui.ReadUserSecret("Enter recovery key: ")
	if err != nil {
		return "", err
	}

	return key, nil
}

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
			if err := devicePathExists(device); err != nil {
				return err
			}

			p, err := ensurePathIsAbolute(dir)
			if err != nil {
				return err
			}

			if err := tpm.MountVolume(device, p, &authRequestor{}); err != nil {
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
			p, err := ensurePathIsAbolute(dir)
			if err != nil {
				return err
			}

			if err := tpm.UnmountVolume(p); err != nil {
				return err
			}

			return nil
		},
	}
}

func newGetLuksKeyFromRecoveryKeyCmd() *cli.Command {
	var hex, escaped bool

	return &cli.Command{
		Name:    "get-luks-key",
		Usage:   "Get LUKS key from recovery key",
		Suggest: true,
		MutuallyExclusiveFlags: []cli.MutuallyExclusiveFlags{
			{
				Flags: [][]cli.Flag{
					{
						&cli.BoolFlag{
							Name:        "hex",
							Usage:       "Output key in hexadecimal format",
							Destination: &hex,
						},
					},
					{
						&cli.BoolFlag{
							Name:        "escaped",
							Usage:       "Output key in escaped string format",
							Destination: &escaped,
						},
					},
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			recoveryKey, err := tui.ReadUserSecret("Enter recovery key: ")
			if err != nil {
				return err
			}

			key, err := tpm.GetLuksKey(ctx, recoveryKey)
			if err != nil {
				return err
			}

			format := "LUKS key (hex): %x\n"
			switch {
			case hex:
				format = "%x\n"
			case escaped:
				format = "%q\n"
			}

			fmt.Printf(format, key)

			return nil
		},
	}
}

// devicePathExists validates that a device path exists in the system.
func devicePathExists(p string) error {
	if p == "" {
		return fmt.Errorf("device path cannot be empty")
	}

	// Check if the device actually exists
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("device %q does not exist", p)
		}
		return fmt.Errorf("failed to check device %q: %v", p, err)
	}

	return nil
}

// ensurePathIsAbolute resolves to an absolute path.
func ensurePathIsAbolute(p string) (string, error) {
	if p == "" {
		return "", fmt.Errorf("directory path cannot be empty")
	}

	if filepath.IsAbs(p) {
		return p, nil
	}

	// Relative path: resolve against current working directory
	r, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not resolve current directory. Please use an absolute path")
	}
	return filepath.Join(r, p), nil
}
