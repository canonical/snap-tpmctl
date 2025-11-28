package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	sm "github.com/egregors/sortedmap"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
)

func newEnumerateCmd() *cli.Command {
	// TODO: add possibility to filter by volume or container role

	return &cli.Command{
		Name:    "list",
		Usage:   "Enumerate all the keyslots",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return enumerate(ctx)
		},
	}
}

func enumerate(ctx context.Context) error {
	c := snapd.NewClient()
	defer c.Close()

	if err := c.LoadAuthFromHome(); err != nil {
		return fmt.Errorf("failed to load auth: %w", err)
	}

	res, err := c.EnumerateKeySlots(ctx)
	if err != nil {
		return err
	}

	if err = displayTable(res); err != nil {
		return err
	}

	return nil
}

func displayTable(data *snapd.SystemVolumesResult) error {
	dashIfEmpty := func(s string) string {
		if strings.TrimSpace(s) == "" {
			return "-"
		}
		return s
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ContainerRole", "Volume", "VolumeName", "Encrypted", "Name", "AuthMode", "PlatformName", "Roles", "Type")

	sortedData := sm.NewFromMap(data.ByContainerRole, func(i, j sm.KV[string, snapd.VolumeInfo]) bool {
		return i.Key < j.Key
	})

	for role, volume := range sortedData.All() {
		keyslots := sm.NewFromMap(volume.KeySlots, func(i, j sm.KV[string, snapd.KeySlotInfo]) bool {
			return i.Key < j.Key
		})

		// TODO: find a better way to do this

		if keyslots.Len() == 0 {
			err := table.Append(
				role,
				volume.Name,
				volume.VolumeName,
				fmt.Sprintf("%v", volume.Encrypted),
				"-",
				"-",
				"-",
				"-",
				"-",
			)
			if err != nil {
				return fmt.Errorf("failed to append table row: %w", err)
			}
		}

		for name, slot := range keyslots.All() {
			err := table.Append(
				role,
				volume.Name,
				volume.VolumeName,
				fmt.Sprintf("%v", volume.Encrypted),
				dashIfEmpty(name),
				dashIfEmpty(slot.AuthMode),
				dashIfEmpty(slot.PlatformName),
				dashIfEmpty(strings.Join(slot.Roles, "+")),
				dashIfEmpty(slot.Type),
			)
			if err != nil {
				return fmt.Errorf("failed to append table row: %w", err)
			}
		}
	}

	if err := table.Render(); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	return nil
}
