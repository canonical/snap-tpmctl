package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func newStatusCmd() *cli.Command {
	return &cli.Command{
		Name:    "status",
		Usage:   "Show TPM status",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.IntArg{
				Name:      "key-id",
				UsageText: "<key-id>",
				Value:     -1,
			},
		},
		Action: status,
	}
}

func status(ctx context.Context, cmd *cli.Command) error {
	// TODO: add validator for key-id
	if cmd.IntArg("key-id") < 0 {
		return cli.Exit("Missing key-id argument", 1)
	}

	fmt.Println("This is my status for key", cmd.IntArg("key-id"))
	slog.Debug("this is my debug log")

	return nil
}
