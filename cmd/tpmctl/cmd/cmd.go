package cmd

import (
	"context"
	"log/slog"

	"github.com/canonical/snap-tpmctl/internal/log"
	"github.com/urfave/cli/v3"
)

/*
  TODO: investigate this
  2025/11/20 11:45:59 ERROR flag needs an argument: --file
  11:46:54 ERROR invalid value "asdf" for argument key-id: strconv.ParseInt: parsing "asdf": invalid syntax
*/

// App is the main application structure.
type App struct {
	args []string
	root cli.Command
}

// New returns a new App.
func New(args []string) App {
	return App{
		args: args,
		root: newRootCmd(),
	}
}

// Run is the main entry point of the app.
func (a App) Run(ctx context.Context) error {
	return a.root.Run(ctx, a.args)
}

func newRootCmd() cli.Command {
	var verbosity int

	// Custom cli version flag
	// TODO: This should be a command itself, not a flag.
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print the version",
	}

	return cli.Command{
		Name:                   "snap-tpmctl",
		Usage:                  "Ubuntu TPM and FDE management tool",
		Version:                "0.1.0",
		UseShortOptionHandling: true,
		EnableShellCompletion:  true,
		Commands: []*cli.Command{
			newAddPINCmd(),
			newAddPassphraseCmd(),
			newCreateKeyCmd(),
			newCheckCmd(),
			newGetLuksKeyFromRecoveryKeyCmd(),
			newListAllCmd(),
			newListPassphraseCmd(),
			newListPinCmd(),
			newListRecoveryKeyCmd(),
			newMountVolumeCmd(),
			newReplacePassphraseCmd(),
			newReplacePinCmd(),
			newRegenerateKeyCmd(),
			newRemovePINCmd(),
			newRemovePassphraseCmd(),
			newStatusCmd(),
			newUnmountVolumeCmd(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbosity",
				Usage:   "Increase verbosity level",
				Aliases: []string{"v"},
				Config: cli.BoolConfig{
					Count: &verbosity,
				},
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			setupLogging(ctx, verbosity)
			return ctx, nil
		},
	}
}

// setupLogging sets up the logging level based on verbosity.
func setupLogging(ctx context.Context, verbosity int) {
	switch verbosity {
	case 0:
		log.SetLoggerLevelInContext(ctx, slog.LevelWarn)
	case 1:
		log.SetLoggerLevelInContext(ctx, slog.LevelInfo)
	default:
		log.SetLoggerLevelInContext(ctx, slog.LevelDebug)
	}
}
