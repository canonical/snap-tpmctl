package cmd_test

import (
	"testing"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	cmdtestutils "github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/matryer/is"
)

func TestRun(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	ctx, logs := testutils.TestLoggerWithBuffer(t)

	app := cmd.New(cmdtestutils.WithArgs("help"))

	err := app.Run(ctx)
	is.NoErr(err)

	is.True(logs.Len() == 0) // No logs printed by default
}

func TestVersion(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	ctx, logs := testutils.TestLoggerWithBuffer(t)

	app := cmd.New(cmdtestutils.WithArgs("version"))

	err := app.Run(ctx)
	is.NoErr(err)

	is.True(logs.Len() == 0) // No logs printed by default
}
