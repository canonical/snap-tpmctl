package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	cmdtestutils "github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/golden"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/creack/pty"
	"github.com/matryer/is"
)

func TestMountVolume(t *testing.T) {
	tests := map[string]struct {
		device string
		dir    string

		recoveryKey       string
		syscall           tpmtestutils.TestSyscall
		deviceInUse       bool
		emptyDeviceError  bool
		emptyDirError     bool
		ttyReadError      bool
		deviceStatError   bool
		alreadyMountedErr bool

		wantErr bool
	}{
		"Success on mounting volume": {},

		"Error out when authRequestor fails":                      {ttyReadError: true, wantErr: true},
		"Error out when mount fails":                              {syscall: tpmtestutils.TestSyscall{WantErr: true}, wantErr: true},
		"Error out when volume is already mounted":                {alreadyMountedErr: true, wantErr: true},
		"Error out when device doesn't exists":                    {deviceStatError: true, wantErr: true},
		"Error out when device is empty":                          {emptyDeviceError: true, wantErr: true},
		"Error out when dir path is empty":                        {emptyDirError: true, wantErr: true},
		"Error out when device is already in use by another tool": {deviceInUse: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			ctx, logs := testutils.TestLoggerWithBuffer(t)

			root := t.TempDir()

			// cryptsetup mock binary
			tpmtestutils.SetupMockBinary(is, root)
			t.Setenv("SNAP", root)

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

			if tc.ttyReadError {
				tty = nil
			}

			var out strings.Builder
			tui := tui.New(tty, &out)

			done := make(chan struct{})
			go func() {
				defer close(done)
				fmt.Fprintln(ptmx, tc.recoveryKey)
			}()

			if tc.emptyDeviceError {
				tc.device = ""
			}

			if tc.emptyDirError {
				tc.dir = ""
			}

			s := tpm.New(
				tpmtestutils.WithRoot(root),
				tpmtestutils.WithSyscall(&tc.syscall),
			)
			app := cmd.New(
				cmdtestutils.WithSnapTPM(s),
				cmdtestutils.WithArgs(command, tc.device, tc.dir),
				cmdtestutils.WithTui(tui),
			)

			err = app.Run(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default
		})
	}
}

func TestUnmountVolume(t *testing.T) {
	tests := map[string]struct {
		dir           string
		syscall       tpmtestutils.TestSyscall
		emptyDirError bool

		wantErr bool
	}{
		"Success on unmounting volume": {},

		"Error out when unmount fails":     {syscall: tpmtestutils.TestSyscall{WantErr: true}, wantErr: true},
		"Error out when dir path is empty": {emptyDirError: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			ctx, logs := testutils.TestLoggerWithBuffer(t)

			root := t.TempDir()

			// cryptsetup mock binary
			tpmtestutils.SetupMockBinary(is, root)
			t.Setenv("SNAP", root)

			command := "unmount-volume"

			if tc.dir == "" {
				tc.dir = "mount-dir"
			}
			tc.dir = filepath.Join(root, tc.dir) // Convert to an absolute path

			var out strings.Builder
			tui := tui.New(nil, &out)

			content := fmt.Sprintf("%s %s ext4 rw 0 0\n", filepath.Join(root, "dev", "mapper", "test"), tc.dir)
			tpmtestutils.SetupProcMount(is, root, content)

			if tc.emptyDirError {
				tc.dir = ""
			}

			s := tpm.New(
				tpmtestutils.WithRoot(root),
				tpmtestutils.WithSyscall(&tc.syscall),
			)
			app := cmd.New(
				cmdtestutils.WithSnapTPM(s),
				cmdtestutils.WithArgs(command, tc.dir),
				cmdtestutils.WithTui(tui),
			)

			err := app.Run(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default
		})
	}
}

func TestGetLuksKeyFromRecoveryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		recoveryKey  string
		hexFlag      bool
		escapedFlag  bool
		ttyReadError bool

		wantErr bool
	}{
		"Success_getting_luks_key":                   {},
		"Success_getting_luks_key_with_hex_flag":     {hexFlag: true},
		"Success_getting_luks_key_with_escaped_flag": {escapedFlag: true},

		"Fail_reading_input":      {ttyReadError: true, wantErr: true},
		"Fail_getting_luks_key":   {recoveryKey: "invalid", wantErr: true},
		"Fail_for_too_many_flags": {wantErr: true, hexFlag: true, escapedFlag: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)
			ctx, logs := testutils.TestLoggerWithBuffer(t)

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

			if tc.ttyReadError {
				tty = nil
			}

			var out strings.Builder
			tui := tui.New(tty, &out)

			done := make(chan struct{})
			go func() {
				defer close(done)
				fmt.Fprintf(ptmx, "%s\n", tc.recoveryKey)
			}()

			app := cmd.New(
				cmdtestutils.WithArgs(args...),
				cmdtestutils.WithTui(tui),
			)

			err = app.Run(ctx)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.True(logs.Len() == 0) // No logs printed by default

			golden.CheckOrUpdate(t, out.String()) // TestGetLuksKeyFromRecoveryKey returns the expected output
		})
	}
}

func TestMain(m *testing.M) {
	if filepath.Base(os.Args[0]) == "systemd-cryptsetup" {
		tpmtestutils.SystemdCryptsetupMock()
		return
	}

	m.Run()
}
