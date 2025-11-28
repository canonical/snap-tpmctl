package cmd

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/log"
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
func (a App) Run() error {
	return a.root.Run(context.Background(), a.args)
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
			newCreateEnterpriseKeyCmd(),
			newCreateKeyCmd(),
			newCheckCmd(),
			newEnumerateCmd(),
			newGetLuksPassphraseCmd(),
			newMountVolumeCmd(),
			newRegenerateEnterpriseKeyCmd(),
			newRegenerateKeyCmd(),
			newStatusCmd(),
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
			setupLogging(verbosity)
			return ctx, nil
		},
	}
}

func setupLogging(level int) {
	switch level {
	case 0:
		log.SetLevel(log.WarnLevel)
	case 1:
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}

	log.SetOutput(os.Stderr)
}
