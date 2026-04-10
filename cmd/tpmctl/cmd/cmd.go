package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/canonical/snap-tpmctl/internal/log"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/urfave/cli/v3"
)

// App is the main application structure.
type App struct {
	option
}

type option struct {
	args []string
	euid int
	tpm  tpm.SnapTPM
	tui  tui.Tui
}

// Option is a functional option for configuring the App.
type Option func(*option)

// New returns a new App.
func New(args ...Option) App {
	o := option{
		args: os.Args,
		euid: os.Geteuid(),
		tpm:  tpm.New(),
		tui:  tui.New(os.Stdin, os.Stdout),
	}
	for _, f := range args {
		f(&o)
	}

	return App{
		option: o,
	}
}

// Run is the main entry point of the app.
func (a App) Run(ctx context.Context) error {
	root := a.newRootCmd()
	return root.Run(ctx, a.args)
}

// isUserRoot returns true if the effective user ID is 0 (root).
func (a App) isUserRoot() bool {
	return a.euid == 0
}

// version is set at build time via ldflags.
var version = "dev"

func (a App) newRootCmd() cli.Command {
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
			a.newAddPINCmd(),
			a.newAddPassphraseCmd(),
			a.newCreateKeyCmd(),
			a.newCheckCmd(),
			a.newGetLuksKeyFromRecoveryKeyCmd(),
			a.newListAllCmd(),
			a.newListPassphraseCmd(),
			a.newListPINCmd(),
			a.newListRecoveryKeyCmd(),
			a.newMountVolumeCmd(),
			a.newReplacePassphraseCmd(),
			a.newReplacePINCmd(),
			a.newRegenerateKeyCmd(),
			a.newRemovePassphraseCmd(),
			a.newRemovePINCmd(),
			a.newStatusCmd(),
			a.newUnmountVolumeCmd(),
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
