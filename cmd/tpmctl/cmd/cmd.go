package cmd

import (
	"context"
	"log/slog"

	"github.com/canonical/snap-tpmctl/internal/log"
	"github.com/urfave/cli/v3"
)

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

// version is set at build time via ldflags.
var version = "dev"

func newRootCmd() cli.Command {
	var verbosity int

	return cli.Command{
		Name:                   "snap-tpmctl",
		Usage:                  "Ubuntu TPM and FDE management tool",
		Version:                version,
		UseShortOptionHandling: true,
		EnableShellCompletion:  true,
		ConfigureShellCompletionCommand: func(cmd *cli.Command) {
			complete := cmd.Action
			cmd.Action = func(ctx context.Context, cmd *cli.Command) error {
				// Patch the output of the bash completion command to remove the "default" file completion on bash.
				if cmd.Args().Len() > 0 && cmd.Args().First() == "bash" {
					cmd.Writer = bashCompletionWriter{Writer: cmd.Writer}
				}

				return complete(ctx, cmd)
			}
		},
		HideVersion: true,
		Commands: []*cli.Command{
			newAddPINCmd(),
			newAddPassphraseCmd(),
			newCreateKeyCmd(),
			newCheckCmd(),
			newGetLuksKeyFromRecoveryKeyCmd(),
			newListAllCmd(),
			newListPassphraseCmd(),
			newListPINCmd(),
			newListRecoveryKeyCmd(),
			newMountVolumeCmd(),
			newReplacePassphraseCmd(),
			newReplacePINCmd(),
			newRegenerateKeyCmd(),
			newRemovePassphraseCmd(),
			newRemovePINCmd(),
			newStatusCmd(),
			newUnmountVolumeCmd(),
			newVersionCmd(),
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
