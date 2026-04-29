// Package tpmtestutils provides helpers for TPM-related tests.
package tpmtestutils

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	_ "unsafe" // Required for go:linkname directives

	"github.com/canonical/snap-tpmctl/internal/snapd"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/matryer/is"
)

func init() {
	testsdetection.MustBeTesting()
}

// WithSnapdClient is an option that configures the TPM to use the provided snapd client.
//
//go:linkname WithSnapdClient github.com/canonical/snap-tpmctl/internal/tpm.withSnapdClient
func WithSnapdClient(snapdClient *snapd.Client) tpm.Option

// WithRoot is an option that configures the TPM to use the provided system mounter.
//
//go:linkname WithRoot github.com/canonical/snap-tpmctl/internal/tpm.withRoot
func WithRoot(r string) tpm.Option

// syscaller abstracts mount and unmount system calls used by SnapTPM.
type syscaller interface {
	Mount(path, target string) error
	Unmount(target string) error
}

// WithSyscall is an option that configures the TPM to use the provided system mounter.
//
//go:linkname WithSyscall github.com/canonical/snap-tpmctl/internal/tpm.withSyscall
func WithSyscall(s syscaller) tpm.Option

// LuksVolumeName converts a directory path into a valid LUKS volume name.
//
//go:linkname LuksVolumeName github.com/canonical/snap-tpmctl/internal/tpm.luksVolumeName
func LuksVolumeName(p string) string

// OneRequestBodyContains checks that at least one request contains all the expected wanted contents.
func OneRequestBodyContains(is *is.I, requests []snapdtestutils.RecordedRequest, wants ...string) {
	is.Helper()

	if len(wants) == 0 {
		panic("Programmer error: RequestsBodyContains checked for nothings as it doesn’t have any wants")
	}

	is.True(slices.ContainsFunc(requests, func(r snapdtestutils.RecordedRequest) bool {
		for _, want := range wants {
			if !strings.Contains(r.Body, want) {
				return false
			}
		}
		return true
	})) // Didn't find all wants in any of the requests.
}

// SystemdCryptsetupMock emulates the systemd-cryptsetup binary behavior for tests.
func SystemdCryptsetupMock() {
	flag.Parse()
	args := flag.Args()

	fmt.Println("Mock systemd-cryptsetup called with args:", args)

	volumeName := args[1]
	if strings.Contains(volumeName, "exit-with-failure") {
		os.Exit(1)
	}

	os.Exit(0)
}

// SetupMockBinary creates a mock systemd-cryptsetup binary in the provided root for tests.
func SetupMockBinary(is *is.I, root string) {
	is.Helper()

	usrBin := filepath.Join(root, "usr", "bin")
	err := os.MkdirAll(usrBin, 0750)
	is.NoErr(err) // Setup: could not create mock binary directory
	path, err := filepath.Abs(os.Args[0])
	is.NoErr(err) // Setup: could not find asbsolute path to self
	err = os.Symlink(path, filepath.Join(usrBin, "systemd-cryptsetup"))
	is.NoErr(err) // Setup: could not create symlink for mock cryptsetup binary
}

// SetupProcMount creates a mock /proc/mounts file with the provided content under root.
func SetupProcMount(is *is.I, root, content string) {
	is.Helper()

	err := os.MkdirAll(filepath.Join(root, "proc"), 0750)
	is.NoErr(err)
	f, err := os.Create(filepath.Join(root, "proc", "mounts"))
	is.NoErr(err)
	defer f.Close()

	_, err = f.WriteString(content)
	is.NoErr(err)
}

// TestSyscall is a test implementation of mount and unmount system calls.
type TestSyscall struct {
	Mounted   bool
	Unmounted bool

	WantErr bool
}

// Mount records a mount call and optionally returns a test error.
func (t *TestSyscall) Mount(path, target string) error {
	if t.WantErr {
		return errors.New("test error")
	}
	t.Mounted = true
	return nil
}

// Unmount records an unmount call and optionally returns a test error.
func (t *TestSyscall) Unmount(target string) error {
	if t.WantErr {
		return errors.New("test error")
	}
	t.Unmounted = true
	return nil
}
