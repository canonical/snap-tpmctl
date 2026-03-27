package cmd_test

import (
	"strings"
	"testing"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	cmdtestutils "github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/matryer/is"
)

func TestRun(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	ctx, logs := testutils.TestLoggerWithBuffer(t)

	var out strings.Builder
	tui := tui.New(nil, &out)

	app := cmd.New(cmdtestutils.WithArgs("help"), cmdtestutils.WithTui(tui))

	err := app.Run(ctx)
	is.NoErr(err)

	is.True(logs.Len() == 0) // No logs printed by default
}

func TestVersion(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	ctx, logs := testutils.TestLoggerWithBuffer(t)

	var out strings.Builder
	tui := tui.New(nil, &out)

	app := cmd.New(cmdtestutils.WithArgs("version"), cmdtestutils.WithTui(tui))

	err := app.Run(ctx)
	is.NoErr(err)

	is.True(logs.Len() == 0) // No logs printed by default
}
