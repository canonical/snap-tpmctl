package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/urfave/cli/v3"
)

func newListAllCmd() *cli.Command {
	var hideHeaders bool

	return &cli.Command{
		Name:    "list-all",
		Usage:   "List all the keyslots with details",
		Suggest: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "no-headers",
				Usage:       "Hide column headers",
				Destination: &hideHeaders,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()

			result, err := c.EnumerateKeySlots(ctx)
			if err != nil {
				return err
			}

			if err := displayAllKeys(os.Stdout, result, hideHeaders); err != nil {
				return err
			}

			return nil
		},
	}
}

func newListPassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:    "list-passphrases",
		Usage:   "List passphrases",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()

			result, err := c.EnumerateKeySlots(ctx)
			if err != nil {
				return err
			}

			data := parseKeySlots(result, (*snapd.KeySlotInfo).IsPassphrase)

			displayKeySlotsFromMap(os.Stdout, "Passphrases", data)

			return nil
		},
	}
}

func newListRecoveryKeyCmd() *cli.Command {
	return &cli.Command{
		Name:    "list-recovery-keys",
		Usage:   "List recovery keys",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()

			result, err := c.EnumerateKeySlots(ctx)
			if err != nil {
				return err
			}

			data := parseKeySlots(result, (*snapd.KeySlotInfo).IsRecoveryKey)

			displayKeySlotsFromMap(os.Stdout, "Recovery Keys", data)

			return nil
		},
	}
}

func newListPinCmd() *cli.Command {
	return &cli.Command{
		Name:    "list-pins",
		Usage:   "List pins",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()

			result, err := c.EnumerateKeySlots(ctx)
			if err != nil {
				return err
			}

			data := parseKeySlots(result, (*snapd.KeySlotInfo).IsPin)

			displayKeySlotsFromMap(os.Stdout, "Pins", data)

			return nil
		},
	}
}

type key struct {
	snapd.KeySlotInfo
	containerRole string
	encrypted     bool
	keySlotName   string
	volume        string
	volumeName    string
}

func getAllKeys(data *snapd.SystemVolumesResult) []key {
	if data == nil {
		return nil
	}

	var volumes []string
	for v := range data.ByContainerRole {
		volumes = append(volumes, v)
	}
	slices.Sort(volumes)

	var allKeys []key
	for _, role := range volumes {
		name := data.ByContainerRole[role].Name
		encrypted := data.ByContainerRole[role].Encrypted
		volumeName := data.ByContainerRole[role].VolumeName

		var keys []key
		for k, v := range data.ByContainerRole[role].KeySlots {
			keys = append(keys, key{
				keySlotName:   k,
				KeySlotInfo:   v,
				containerRole: role,
				volume:        name,
				encrypted:     encrypted,
				volumeName:    volumeName,
			})
		}

		slices.SortFunc(keys, func(i, j key) int {
			return strings.Compare(i.keySlotName, j.keySlotName)
		})

		allKeys = append(allKeys, keys...)
	}

	return allKeys
}

func displayAllKeys(w io.Writer, data *snapd.SystemVolumesResult, hideHeaders bool) error {
	keys := getAllKeys(data)

	rows := [][]string{}
	dashIfEmpty := func(s string) string {
		if strings.TrimSpace(s) == "" {
			return "-"
		}
		return s
	}

	for _, k := range keys {
		rows = append(rows, []string{
			k.containerRole,
			k.volume,
			k.volumeName,
			fmt.Sprintf("%v", k.encrypted),
			dashIfEmpty(k.keySlotName),
			dashIfEmpty(k.AuthMode),
			dashIfEmpty(k.PlatformName),
			dashIfEmpty(strings.Join(k.Roles, "+")),
			dashIfEmpty(k.Type),
		})
	}

	headers := []string{"ContainerRole", "Volume", "VolumeName", "Encrypted", "KeyslotName", "AuthMode", "PlatformName", "Roles", "Type"}

	if err := tui.DisplayTable(w, headers, rows, hideHeaders); err != nil {
		return err
	}

	return nil
}

func parseKeySlots(data *snapd.SystemVolumesResult, filter func(*snapd.KeySlotInfo) bool) []string {
	keys := getAllKeys(data)

	var result []string
	for _, k := range keys {
		if !filter(&k.KeySlotInfo) {
			continue
		}

		// Deduplicate entries
		if slices.Contains(result, k.keySlotName) {
			continue
		}

		result = append(result, k.keySlotName)
	}

	return result
}

func displayKeySlotsFromMap(w io.Writer, title string, entries []string) {
	fmt.Fprintf(w, "%s:\n", title)

	for _, e := range entries {
		fmt.Fprintf(w, "* %s\n", e)
	}
}
