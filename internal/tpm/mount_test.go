package tpm_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/matryer/is"
	"github.com/snapcore/secboot"
)

func TestMountVolume(t *testing.T) {
	tests := map[string]struct {
		device            string
		target            string
		syscall           tpmtestutils.TestSyscall
		authRequestor     authRequestor
		targetExists      bool
		mkdirErr          bool
		alreadyMountedErr bool

		wantMounted   bool
		wantRequested bool

		wantErr bool
	}{
		"Success on mounting volume": {wantRequested: true, wantMounted: true},
		"Success when target already exists": {
			target:        "existing-mount-dir",
			targetExists:  true,
			wantRequested: true,
			wantMounted:   true,
		},

		"Error out when unable to crate directory": {mkdirErr: true, wantErr: true},
		"Error out when authRequestor fails":       {authRequestor: authRequestor{wantErr: true}, wantErr: true},
		"Error out when unable to mount volume":    {syscall: tpmtestutils.TestSyscall{WantErr: true}, wantRequested: true, wantErr: true},
		"Error out when volume is already mounted": {alreadyMountedErr: true, wantErr: true},
		"Error out when systemd-cryptsetup fails":  {device: "exit-with-failure", wantRequested: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			root := t.TempDir()

			// cryptsetup mock binary
			tpmtestutils.SetupMockBinary(is, root)
			t.Setenv("SNAP", root)

			if tc.device == "" {
				tc.device = "test-device"
			}
			tc.device = filepath.Join(root, tc.device) // Convert to an absolute path

			if tc.target == "" {
				tc.target = "mount-dir"
			}
			tc.target = filepath.Join(root, tc.target) // Convert to an absolute path

			content := ""
			if tc.alreadyMountedErr {
				mapper := filepath.Join(root, "dev/mapper", tpmtestutils.LuksVolumeName(tc.device))
				content = fmt.Sprintf("%s %s ext4 rw 0 0\n", mapper, tc.target)
			}
			tpmtestutils.SetupProcMount(is, root, content)

			if tc.targetExists {
				err := os.MkdirAll(tc.target, 0750)
				is.NoErr(err) // Setup: target directory should exist before mounting
			}

			// In order to test the `MkdirAll` failure, we need create a target file with the supposed target folder name.
			if tc.mkdirErr {
				f, err := os.Create(tc.target)
				is.NoErr(err) // Setup: target should exist before mounting as a file
				defer f.Close()
			}

			s := tpm.New(
				tpmtestutils.WithRoot(root),
				tpmtestutils.WithSyscall(&tc.syscall),
			)

			err := s.Mount(ctx, tc.device, tc.target, &tc.authRequestor)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.Equal(tc.authRequestor.requested, tc.wantRequested) // the recovery key is asked as expected
			is.Equal(tc.syscall.Mounted, tc.wantMounted)           // the volume is mounted as expected
		})
	}
}

func TestUnmountVolume(t *testing.T) {
	tests := map[string]struct {
		target  string
		mapper  string
		syscall tpmtestutils.TestSyscall

		wantUnmounted bool

		wantErr      bool
		wantRmdirErr bool
	}{
		"Success on unmounting volume": {wantUnmounted: true},

		"Error out when unable to remove directory":   {wantRmdirErr: true, wantErr: true},
		"Error out when unable determine device path": {target: "not-existing-target", wantErr: true},
		"Error out when unable to unmount volume":     {syscall: tpmtestutils.TestSyscall{WantErr: true}, wantErr: true},
		"Error out when systemd-cryptsetup fails":     {mapper: "exit-with-failure", wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			root := t.TempDir()

			// cryptsetup mock binary
			tpmtestutils.SetupMockBinary(is, root)
			t.Setenv("SNAP", root)

			if tc.mapper == "" {
				tc.mapper = "test-device"
			}
			tc.mapper = filepath.Join(root, "dev", "mapper", tc.mapper) // Convert to an absolute path

			target := "mount-dir"
			if tc.target == "" {
				tc.target = target
			}
			tc.target = filepath.Join(root, tc.target) // Convert to an absolute path

			content := fmt.Sprintf("%s %s ext4 rw 0 0\n", tc.mapper, filepath.Join(root, target))
			tpmtestutils.SetupProcMount(is, root, content)

			// In order to test the `RemoveAll` failure, we need to set restrictive permissions for the target's parent folder.
			if tc.wantRmdirErr {
				err := os.MkdirAll(tc.target, 0750)
				is.NoErr(err)

				//nolint:gosec // test-only permissions, non-sensitive temp path
				err = os.Chmod(filepath.Dir(tc.target), 0555)
				is.NoErr(err)

				defer func() {
					//nolint:gosec // test-only permissions, non-sensitive temp path
					err := os.Chmod(filepath.Dir(tc.target), 0750)
					is.NoErr(err)
				}()
			}

			s := tpm.New(
				tpmtestutils.WithRoot(root),
				tpmtestutils.WithSyscall(&tc.syscall),
			)

			err := s.Unmount(ctx, tc.target)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.Equal(tc.syscall.Unmounted, tc.wantUnmounted) // the volume is unmounted as expected
		})
	}
}

func TestGetMapperFromMount(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		target  string
		mapper  string
		fileErr bool
		readErr bool

		wantErr bool
	}{
		"Success on getting mapper":                 {},
		"Success with mapper found in /proc/mounts": {target: "wrong-target"},

		"Fail to open /proc/mounts": {fileErr: true, wantErr: true},
		"Fail to read /proc/mounts": {readErr: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			root := t.TempDir()

			if tc.mapper == "" && tc.target == "" {
				tc.mapper = filepath.Join(root, "dev", "mapper", "test-device")
			}

			if tc.target == "" {
				tc.target = "mount-dir"
			}
			tc.target = filepath.Join(root, tc.target) // Convert to an absolute path

			content := fmt.Sprintf("%s %s ext4 rw 0 0\n", tc.mapper, tc.target)
			if tc.readErr {
				// Scanner default max token: 64K. This will return a read error
				content = strings.Repeat("a", 70*1024) + "\n"
			}

			tpmtestutils.SetupProcMount(is, root, content)

			if tc.fileErr {
				err := os.Remove(filepath.Join(root, "proc", "mounts"))
				is.NoErr(err) // Setup: /proc/mounts should be deleted for a file error
			}

			s := tpm.New(tpmtestutils.WithRoot(root))

			m, err := tpm.GetMapperFromMount(s, tc.target)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.Equal(m, tc.mapper) // the device mapper is the expected one
		})
	}
}

func TestGetMountFromMapper(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		mapper  string
		mount   string
		fileErr bool
		readErr bool

		wantErr bool
	}{
		"Success on getting mount":                 {},
		"Success with mount found in /proc/mounts": {mapper: "wrong-mapper"},

		"Fail to open /proc/mounts": {fileErr: true, wantErr: true},
		"Fail to read /proc/mounts": {readErr: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			root := t.TempDir()

			if tc.mount == "" && tc.mapper == "" {
				tc.mount = "mount-dir"
			}

			mapper := filepath.Join(root, "dev", "mapper", "test-device") // Convert to an absolute path
			if tc.mapper == "" {
				tc.mapper = mapper
			}

			content := fmt.Sprintf("%s %s ext4 rw 0 0\n", mapper, tc.mount)
			if tc.readErr {
				// Scanner default max token: 64K. This will return a read error
				content = strings.Repeat("a", 70*1024) + "\n"
			}

			tpmtestutils.SetupProcMount(is, root, content)

			if tc.fileErr {
				err := os.Remove(filepath.Join(root, "proc", "mounts"))
				is.NoErr(err) // Setup: /proc/mounts should be deleted for a file error
			}

			s := tpm.New(tpmtestutils.WithRoot(root))

			m, err := tpm.GetMountFromMapper(s, tc.mapper)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.Equal(m, tc.mount) // the mount path is the expected one
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

type authRequestor struct {
	requested bool

	wantErr bool
}

func (r *authRequestor) RequestUserCredential(ctx context.Context, name, path string, authTypes secboot.UserAuthType) (string, secboot.UserAuthType, error) {
	if r.wantErr {
		return "", 0, errors.New("test error")
	}
	r.requested = true
	return "22003-18216-51619-31723-49692-17125-14174-57839", secboot.UserAuthTypeRecoveryKey, nil
}

func (r *authRequestor) NotifyUserAuthResult(ctx context.Context, result secboot.UserAuthResult, authTypes, exhaustedAuthTypes secboot.UserAuthType) error {
	return nil
}
