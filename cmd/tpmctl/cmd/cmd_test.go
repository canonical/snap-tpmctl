package cmd_test

import (
	"testing"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/matryer/is"
)

// TODO: add tests to every subcommands.

func TestRun(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	ctx, logs := testutils.TestLoggerWithBuffer(t)

	app := cmd.New([]string{"--help"})

	err := app.Run(ctx)
	is.NoErr(err)

	is.True(logs.Len() == 0) // No logs printed by default
}
