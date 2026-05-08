package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/creack/pty"
	"github.com/matryer/is"
)

func TestMountVolume(t *testing.T) {
	tests := map[string]struct {
		device      string
		dir         string
		recoveryKey string

		mountErr          bool
		deviceInUse       bool
		emptyDeviceError  bool
		emptyDirError     bool
		deviceStatError   bool
		alreadyMountedErr bool

		wantErr bool
	}{
		"Success on mounting volume": {},

		"Error out when mount fails":               {mountErr: true, wantErr: true},
		"Error out when volume is already mounted": {alreadyMountedErr: true, wantErr: true},
		"Error out when device doesn't exists":     {deviceStatError: true, wantErr: true},
		"Error out when device is empty":           {emptyDeviceError: true, wantErr: true},
		"Error out when dir path is empty":         {emptyDirError: true, wantErr: true}, "Error out when device is already in use by another tool": {deviceInUse: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)

			root := t.TempDir()

			// cryptsetup mock binary
			tpmtestutils.SetupMockBinary(is, root)

			command := "mount-volume"

			if tc.device == "" {
				tc.device = "test-device"
			}
			tc.device = filepath.Join(root, tc.device) // Convert to an absolute path

			if !tc.deviceStatError {
				f, err := os.Create(tc.device)
				is.NoErr(err) // Setup: device should exist before mounting
				defer f.Close()
			}

			if tc.dir == "" {
				tc.dir = "mount-dir"
			}
			tc.dir = filepath.Join(root, tc.dir) // Convert to an absolute path

			if tc.recoveryKey == "" {
				tc.recoveryKey = "11272-47509-28031-54818-41671-38673-11053-06376"
			}

			content := ""
			if tc.alreadyMountedErr {
				mapper := tpmtestutils.LuksVolumeName(tc.device)
				content = fmt.Sprintf("%s %s ext4 rw 0 0\n", filepath.Join(root, "dev", "mapper", mapper), tc.dir)
			}
			tpmtestutils.SetupProcMount(is, root, content)
			tpmtestutils.SetupSysClassBlock(is, root, tc.device, tc.deviceInUse)

			ptmx, tty, err := pty.Open()
			is.NoErr(err)
			defer ptmx.Close()
			defer tty.Close()

			go func() {
				fmt.Fprintln(ptmx, tc.recoveryKey)
			}()

			if tc.emptyDeviceError {
				tc.device = ""
			}

			if tc.emptyDirError {
				tc.dir = ""
			}

			//nolint:gosec // The test intentionally executes the binary built in TestMain.
			cmd := exec.Command(cmdPath, command, tc.device, tc.dir)
			cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsRoot(), testutils.WithSyscallErr(tc.mountErr))
			cmd.Stdin = tty

			// no need for checking output here
			_, err = cmd.CombinedOutput()
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestUnmountVolume(t *testing.T) {
	tests := map[string]struct {
		dir string

		unmountErr    bool
		emptyDirError bool

		wantErr bool
	}{
		"Success on unmounting volume": {},

		"Error out when unmount fails":     {unmountErr: true, wantErr: true},
		"Error out when dir path is empty": {emptyDirError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)

			root := t.TempDir()

			// cryptsetup mock binary
			tpmtestutils.SetupMockBinary(is, root)

			command := "unmount-volume"

			if tc.dir == "" {
				tc.dir = "mount-dir"
			}
			tc.dir = filepath.Join(root, tc.dir) // Convert to an absolute path

			content := fmt.Sprintf("%s %s ext4 rw 0 0\n", filepath.Join(root, "dev", "mapper", "test"), tc.dir)
			tpmtestutils.SetupProcMount(is, root, content)

			if tc.emptyDirError {
				tc.dir = ""
			}

			//nolint:gosec // The test intentionally executes the binary built in TestMain.
			cmd := exec.Command(cmdPath, command, tc.dir)
			cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsRoot(), testutils.WithSyscallErr(tc.unmountErr))

			// no need for checking output here
			_, err := cmd.CombinedOutput()
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}
		})
	}
}

func TestGetLuksKeyFromRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKey string
		hexFlag     bool
		escapedFlag bool

		wantErr bool
	}{
		"Success_getting_luks_key":                   {},
		"Success_getting_luks_key_with_hex_flag":     {hexFlag: true},
		"Success_getting_luks_key_with_escaped_flag": {escapedFlag: true},

		"Fail_getting_luks_key":   {recoveryKey: "invalid", wantErr: true},
		"Fail_for_too_many_flags": {wantErr: true, hexFlag: true, escapedFlag: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			root := t.TempDir()

			command := "get-luks-key"

			args := []string{command}
			if tc.hexFlag {
				args = append(args, "--hex")
			}
			if tc.escapedFlag {
				args = append(args, "--escaped")
			}

			if tc.recoveryKey == "" {
				tc.recoveryKey = "11272-47509-28031-54818-41671-38673-11053-06376"
			}

			ptmx, tty, err := pty.Open()
			is.NoErr(err)
			defer ptmx.Close()
			defer tty.Close()

			go func() {
				fmt.Fprintln(ptmx, tc.recoveryKey)
			}()

			cmd := exec.Command(cmdPath, args...)
			cmd.Env = append(cmd.Env, testutils.WithRootDir(root), testutils.WithUserAsNonRoot())
			cmd.Stdin = tty

			out, err := cmd.CombinedOutput()
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			golden.CheckOrUpdate(t, out) // TestGetLuksKeyFromRecoveryKey returns the expected output
		})
	}
}
