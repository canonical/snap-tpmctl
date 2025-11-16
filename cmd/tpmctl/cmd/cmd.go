package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

// TODO: make logging with a more human output. (look at authd for slog)
// TODO: add tests for main.go
// 		Look for table testing
// 		Look at parallel testing

// App is the main application structure.
type App struct {
	root cli.Command
}

// New returns a new App.
func New() App {
	return App{
		root: newRootCmd(),
	}
}

// Run is the main entry point of the app.
func (a App) Run() error {
	return a.root.Run(context.Background(), os.Args)
}

func newRootCmd() cli.Command {
	var verbosity int

	// Custom cli version flag
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
			newStatusCmd(),
			newCreateKeyCmd(),
			newCreateEnterpriseKeyCmd(),
			newRegenerateKeyCmd(),
			newRegenerateEnterpriseKeyCmd(),
			newMountVolumeCmd(),
			newGetLuksPassphraseCmd(),
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
			// &cli.StringFlag{
			// 	Name:  "tpm-path",
			// 	Usage: "TPM device path",
			// 	Value: "/dev/tpm0",
			// },
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			println(verbosity)
			setupLogging(verbosity)
			return ctx, nil
		},
	}
}

func setupLogging(level int) {
	l := slog.LevelWarn
	switch level {
	case 0:
	case 2:
		l = slog.LevelInfo
	default:
		l = slog.LevelDebug
	}

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: l,
	})
	slog.SetDefault(slog.New(h))
}
