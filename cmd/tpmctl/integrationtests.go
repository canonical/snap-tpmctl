//go:build integrationtests

package main

import (
	"os"
	"strconv"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	cmdtestutils "github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd/testutils"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
)

func init() {
	testsdetection.MustBeTesting()

	euid, err := strconv.Atoi(os.Getenv("SNAP_TPMCTL_INTEGRATION_TESTS_ADMIN_EUID"))
	if err != nil {
		panic("SNAP_TPMCTL_INTEGRATION_TESTS_ADMIN_EUID must be a valid integer")
	}

	root := os.Getenv("SNAP_TPMCTL_INTEGRATION_TESTS_ROOT_DIR")
	if root == "" {
		panic("SNAP_TPMCTL_INTEGRATION_TESTS_ROOT_DIR must be set")
	}

	c := snapdtestutils.NewMockSnapdIntegrationServer(root)
	s := tpm.New(tpmtestutils.WithSnapdClient(c.Client))

	mainApp = cmd.New(cmdtestutils.WithSnapTPM(s), cmdtestutils.WithEuid(euid))
}
