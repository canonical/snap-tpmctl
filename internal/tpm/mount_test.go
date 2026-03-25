package tpm_test

import (
	"context"
	"errors"
	"flag"
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
		device        string
		target        string
		syscall       testSyscall
		authRequestor authRequestor
		targetExists  bool

		wantMounted   bool
		wantRequested bool

		wantErr      bool
		wantMkdirErr bool
	}{
		"Success on mounting volume": {wantRequested: true, wantMounted: true},
		"Success when target already exists": {
			target:        "existing-mount-dir",
			targetExists:  true,
			wantRequested: true,
			wantMounted:   true,
		},

		"Error out when unable to crate directory": {wantMkdirErr: true, wantErr: true},
		"Error out when authRequestor fails":       {authRequestor: authRequestor{wantErr: true}, wantErr: true},
		"Error out when unable to mount volume":    {syscall: testSyscall{wantErr: true}, wantRequested: true, wantErr: true},
		"Error out when systemd-cryptsetup fails":  {device: "exit-with-failure", wantRequested: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			root := t.TempDir()

			// cryptsetup mock binary
			setupMockBinary(is, root)
			t.Setenv("SNAP", root)

			if tc.device == "" {
				tc.device = "test-device"
			}
			tc.device = filepath.Join(root, tc.device) // Convert to an absolute path

			if tc.target == "" {
				tc.target = "mount-dir"
			}
			tc.target = filepath.Join(root, tc.target) // Convert to an absolute path

			if tc.targetExists {
				err := os.MkdirAll(tc.target, 0750)
				is.NoErr(err) // Setup: target directory should exist before mounting
			}

			// In order to test the `MkdirAll` failure, we need create a target file with the supposed target folder name.
			if tc.wantMkdirErr {
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
			is.Equal(tc.syscall.mounted, tc.wantMounted)           // the volume is mounted as expected
		})
	}
}

func TestUnmountVolume(t *testing.T) {
	tests := map[string]struct {
		target  string
		mapper  string
		syscall testSyscall

		wantUnmounted bool

		wantErr      bool
		wantRmdirErr bool
	}{
		"Success on unmounting volume": {wantUnmounted: true},

		"Error out when unable to remove directory":   {wantRmdirErr: true, wantErr: true},
		"Error out when unable determine device path": {target: "not-existing-target", wantErr: true},
		"Error out when unable to unmount volume":     {syscall: testSyscall{wantErr: true}, wantErr: true},
		"Error out when systemd-cryptsetup fails":     {mapper: "exit-with-failure", wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			root := t.TempDir()

			// cryptsetup mock binary
			setupMockBinary(is, root)
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
			setupProcMount(is, root, content)

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

			is.Equal(tc.syscall.unmounted, tc.wantUnmounted) // the volume is unmounted as expected
		})
	}
}

func TestGetMapperFromMount(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		target string

		wantFileErr bool
		wantReadErr bool
		wantErr     bool
	}{
		"Success on getting mapper": {},

		"Fail to find mapper /proc/mounts": {target: "wrong-target", wantErr: true},
		"Fail to open /proc/mounts":        {wantFileErr: true, wantErr: true},
		"Fail to read /proc/mounts":        {wantReadErr: true, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			root := t.TempDir()

			mapper := filepath.Join(root, "dev", "mapper", "test-device")
			target := "mount-dir"

			if tc.target == "" {
				tc.target = target
			}
			tc.target = filepath.Join(root, tc.target) // Convert to an absolute path

			content := fmt.Sprintf("%s %s ext4 rw 0 0\n", mapper, filepath.Join(root, target))
			if tc.wantReadErr {
				// Scanner default max token: 64K. This will return a read error
				content = strings.Repeat("a", 70*1024) + "\n"
			}

			setupProcMount(is, root, content)

			if tc.wantFileErr {
				err := os.Remove(filepath.Join(root, "proc", "mounts"))
				is.NoErr(err) // Setup: /proc/mounts should be deleted for a file error
			}

			s := tpm.New(tpmtestutils.WithRoot(root))

			m, err := tpm.GetMapperFromMount(s, tc.target)
			if testutils.CheckError(is, err, tc.wantErr) {
				return
			}

			is.Equal(m, mapper) // the device mapper is the expected one
		})
	}
}

func TestMain(m *testing.M) {
	if filepath.Base(os.Args[0]) == "systemd-cryptsetup" {
		systemdCryptsetupMock()
		return
	}

	m.Run()
}

func systemdCryptsetupMock() {
	flag.Parse()
	args := flag.Args()

	fmt.Println("Mock systemd-cryptsetup called with args:", args)

	volumeName := args[1]
	if strings.Contains(volumeName, "exit-with-failure") {
		os.Exit(1)
	}

	os.Exit(0)
}

func setupMockBinary(is *is.I, root string) {
	is.Helper()

	usrBin := filepath.Join(root, "usr", "bin")
	err := os.MkdirAll(usrBin, 0750)
	is.NoErr(err) // Setup: could not create mock binary directory
	path, err := filepath.Abs(os.Args[0])
	is.NoErr(err) // Setup: could not find asbsolute path to self
	err = os.Symlink(path, filepath.Join(usrBin, "systemd-cryptsetup"))
	is.NoErr(err) // Setup: could not create symlink for mock cryptsetup binary
}

func setupProcMount(is *is.I, root, content string) {
	is.Helper()

	err := os.MkdirAll(filepath.Join(root, "proc"), 0750)
	is.NoErr(err)
	f, err := os.Create(filepath.Join(root, "proc", "mounts"))
	is.NoErr(err)
	defer f.Close()

	_, err = f.WriteString(content)
	is.NoErr(err)
}

type authRequestor struct {
	requested bool

	wantErr bool
}

func (r *authRequestor) RequestUserCredential(ctx context.Context, name, path string, authTypes secboot.UserAuthType) (string, error) {
	if r.wantErr {
		return "", errors.New("test error")
	}
	r.requested = true
	return "22003-18216-51619-31723-49692-17125-14174-57839", nil
}

type testSyscall struct {
	mounted   bool
	unmounted bool

	wantErr bool
}

func (t *testSyscall) Mount(path, target string) error {
	if t.wantErr {
		return errors.New("test error")
	}
	t.mounted = true
	return nil
}

func (t *testSyscall) Unmount(target string) error {
	if t.wantErr {
		return errors.New("test error")
	}
	t.unmounted = true
	return nil
}
