package cmd

import (
	"context"
	"fmt"
	"snap-tpmctl/internal/snapd"

	"github.com/urfave/cli/v3"
)

func newStatusCmd() *cli.Command {
	return &cli.Command{
		Name:    "status",
		Usage:   "Show TPM status",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return status(ctx)
		},
	}
}

func status(ctx context.Context) error {
	// TODO: implement the status API when lands on snapd

	fmt.Println("This is my status for the system")

	c := snapd.NewClient()
	defer c.Close()

	res, err := c.EnumerateKeySlots(ctx)
	if err != nil {
		return err
	}

	for _, volume := range res.ByContainerRole {
		// Look for text template in Go
		// https://pkg.go.dev/text/template

		fmt.Printf("Volume: %s\n", volume.Name)
		fmt.Printf("  Encrypted: %v\n", volume.Encrypted)
		fmt.Printf("  VolumeName: %v\n", volume.VolumeName)

		if len(volume.KeySlots) > 0 {
			fmt.Println("  KeySlots:")
		}

		for _, slot := range volume.KeySlots {
			fmt.Printf("    AuthMode: %v\n", slot.AuthMode)
			fmt.Printf("    PlatformName: %v\n", slot.PlatformName)
			fmt.Printf("    Roles: %v\n", slot.Roles)
			fmt.Printf("    Type: %v\n", slot.Type)
		}
		fmt.Println()
	}

	return nil
}
