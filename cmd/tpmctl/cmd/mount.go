package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func newMountVolumeCmd() *cli.Command {
	return &cli.Command{
		Name:    "mount-volume",
		Usage:   "Unlock and mount a LUKS encrypted volume",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "mount-point",
				UsageText: "<mount-point>",
			},
		},
		Action: mountVolume,
	}
}

func mountVolume(ctx context.Context, cmd *cli.Command) error {
	// TODO: add validator for mount-point [string]

	if cmd.StringArg("mount-point") == "" {
		return cli.Exit("Missing mount-point argument", 1)
	}

	fmt.Println("Mount volume in", cmd.StringArg("mount-point"))
	return nil
}

func newGetLuksPassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:    "get-luks-passphrase",
		Usage:   "Get LUKS passphrase from recovery key",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.IntArg{
				Name:      "key-id",
				UsageText: "<key-id>",
				Value:     -1,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Usage: "Write passphrase to file",
			},
		},
		Action: getLuksPassphrase,
	}
}

func getLuksPassphrase(ctx context.Context, cmd *cli.Command) error {
	// TODO: add validator for key-id [int]

	if cmd.IntArg("key-id") < 0 {
		// TODO… return an error instead
		return cli.Exit("Missing key-id argument", 1)
	}

	// TODO: add validator for f [string]

	msg := "print to stdout"
	if f := cmd.String("file"); f != "" {
		msg = fmt.Sprintf("print to file: %s", f)
	}
	fmt.Println(msg)

	fmt.Println("Get LUKS passphrase for key", cmd.IntArg("key-id"))

	return nil
}
