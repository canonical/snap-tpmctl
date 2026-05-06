//go:build integrationtests

package main

import (
	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	cmdtestutils "github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd/testutils"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
)

func init() {
	testsdetection.MustBeTesting()

	euid := testutils.GetEuidEnv()
	root := testutils.GetRootDirEnv()

	c := snapdtestutils.NewMockSnapdServerWithPath(root)
	s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

	mainApp = cmd.New(cmdtestutils.WithSnapTPM(s), cmdtestutils.WithEuid(euid))
}
