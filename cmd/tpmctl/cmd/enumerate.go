package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	sm "github.com/egregors/sortedmap"
	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
	"snap-tpmctl/internal/tui"
)

func newListDetailCmd() *cli.Command {
	var hideHeaders bool

	return &cli.Command{
		Name:    "list-details",
		Usage:   "List all the keyslots with details",
		Suggest: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "no-headers",
				Usage:       "Show column headers",
				Destination: &hideHeaders,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			result, err := c.EnumerateKeySlots(ctx)
			if err != nil {
				return err
			}

			if err := displayKeysWithDetails(os.Stdout, result, hideHeaders); err != nil {
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
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			result, err := c.EnumerateKeySlots(ctx)
			if err != nil {
				return err
			}

			data := parseKeySlots(result, "passphrase")

			if err := displayKeySlotsFromMap(os.Stdout, "Passphrases", data); err != nil {
				return err
			}

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
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			result, err := c.EnumerateKeySlots(ctx)
			if err != nil {
				return err
			}

			data := parseKeySlots(result, "recovery")

			if err := displayKeySlotsFromMap(os.Stdout, "Recovery Keys", data); err != nil {
				return err
			}

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
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			result, err := c.EnumerateKeySlots(ctx)
			if err != nil {
				return err
			}

			data := parseKeySlots(result, "pin")

			if err := displayKeySlotsFromMap(os.Stdout, "Recovery Keys", data); err != nil {
				return err
			}

			return nil
		},
	}
}

func displayKeysWithDetails(w io.Writer, data *snapd.SystemVolumesResult, hideHeaders bool) error {
	if data == nil {
		return nil
	}

	sortedData := sm.NewFromMap(data.ByContainerRole, func(i, j sm.KV[string, snapd.VolumeInfo]) bool {
		return i.Key < j.Key
	})

	rows := [][]string{}
	dashIfEmpty := func(s string) string {
		if strings.TrimSpace(s) == "" {
			return "-"
		}
		return s
	}

	for role, volume := range sortedData.All() {
		keyslots := sm.NewFromMap(volume.KeySlots, func(i, j sm.KV[string, snapd.KeySlotInfo]) bool {
			return i.Key < j.Key
		})

		if keyslots.Len() == 0 {
			continue
		}

		for name, slot := range keyslots.All() {
			rows = append(rows, []string{
				role,
				volume.Name,
				volume.VolumeName,
				fmt.Sprintf("%v", volume.Encrypted),
				dashIfEmpty(name),
				dashIfEmpty(slot.AuthMode),
				dashIfEmpty(slot.PlatformName),
				dashIfEmpty(strings.Join(slot.Roles, "+")),
				dashIfEmpty(slot.Type),
			})
		}
	}

	headers := []string{"ContainerRole", "Volume", "VolumeName", "Encrypted", "Name", "AuthMode", "PlatformName", "Roles", "Type"}

	if err := tui.DisplayTable(w, headers, rows, hideHeaders); err != nil {
		return err
	}

	return nil
}

func parseKeySlots(data *snapd.SystemVolumesResult, keyType string) []string {
	var result []string

	if data == nil {
		return result
	}

	for _, volume := range data.ByContainerRole {
		for name, slot := range volume.KeySlots {
			if keyType == "recovery" && slot.Type == "recovery" {
				result = append(result, name)
			} else if slot.AuthMode == keyType {
				result = append(result, name)
			}
		}
	}

	// Sort and deduplicate entries
	slices.Sort(result)
	result = slices.Compact(result)

	return result
}

func displayKeySlotsFromMap(w io.Writer, title string, entries []string) error {
	if _, err := fmt.Fprintf(w, "%s:\n", title); err != nil {
		return err
	}

	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "* none")
		return err
	}

	for _, e := range entries {
		if _, err := fmt.Fprintf(w, "* %s\n", e); err != nil {
			return err
		}
	}

	return nil
}
