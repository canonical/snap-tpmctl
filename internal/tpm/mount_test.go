package tpm_test

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	tpmtestutils "github.com/canonical/snap-tpmctl/internal/tpm/testutils"
	"github.com/matryer/is"
	"github.com/snapcore/secboot"
)

func TestMountVolume(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		device string
		target string

		activator  testActivator
		filesystem testFileSystem
		mounter    testMounter

		wantActivated bool
		wantMounted   bool

		wantErr bool
	}{
		"Success on mounting volume": {
			device:        "/dev/test",
			target:        "/media/vol",
			activator:     testActivator{},
			mounter:       testMounter{},
			filesystem:    testFileSystem{},
			wantActivated: true,
			wantMounted:   true,
		},
		"Success on mounting already active volume": {
			device:    "/dev/test",
			target:    "/media/vol",
			activator: testActivator{},
			mounter:   testMounter{},
			filesystem: testFileSystem{
				MapFS: fstest.MapFS{
					"dev/mapper/dev-test": &fstest.MapFile{},
				},
			},
			wantMounted: true,
		},

		"Fail to create directory": {
			device:    "/dev/test",
			target:    "/media/vol",
			activator: testActivator{},
			mounter:   testMounter{},
			filesystem: testFileSystem{
				wantErr: true,
			},
			wantErr: true,
		},
		"Fail to activate volume": {
			device:     "/dev/test",
			target:     "/media/vol",
			activator:  testActivator{wantErr: true},
			mounter:    testMounter{},
			filesystem: testFileSystem{},
			wantErr:    true,
		},
		"Fail to mount volume": {
			device:        "/dev/test",
			target:        "/media/vol",
			activator:     testActivator{},
			mounter:       testMounter{wantErr: true},
			filesystem:    testFileSystem{},
			wantActivated: true,
			wantErr:       true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			m := tpm.NewMount(
				tpmtestutils.WithActivator(&tc.activator),
				tpmtestutils.WithMounter(&tc.mounter),
				tpmtestutils.WithFileSystem(tc.filesystem),
			)

			err := m.MountVolume(ctx, tc.device, tc.target)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.Equal(tc.activator.activated, tc.wantActivated) // the volume is activated as expected
			is.Equal(tc.mounter.mounted, tc.wantMounted)       // the volume is mounted as expected
		})
	}
}

func TestUnmountVolume(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		target string

		activator  testActivator
		filesystem testFileSystem
		mounter    testMounter

		wantDectivated bool
		wantUnmounted  bool

		wantErr bool
	}{
		"Success on deactivating volume": {
			target:    "/media/vol",
			activator: testActivator{},
			mounter:   testMounter{},
			filesystem: testFileSystem{
				MapFS: fstest.MapFS{
					"proc/mounts": &fstest.MapFile{
						Data: []byte("/dev/mapper/dev-test /media/vol ext4 rw 0 0\n"),
					},
				},
			},
			wantDectivated: true,
			wantUnmounted:  true,
		},

		"Fail to get device /proc/mounts": {
			target:     "/media/vol",
			activator:  testActivator{},
			mounter:    testMounter{},
			filesystem: testFileSystem{},
			wantErr:    true,
		},
		"Fail to read /proc/mounts": {
			target:    "/media/vol",
			activator: testActivator{},
			mounter:   testMounter{},
			filesystem: testFileSystem{
				readErr: true,
			},
			wantErr: true,
		},
		"Fail to find device from /proc/mounts": {
			target:    "/media/vol",
			activator: testActivator{},
			mounter:   testMounter{},
			filesystem: testFileSystem{
				MapFS: fstest.MapFS{
					"proc/mounts": &fstest.MapFile{
						Data: []byte("/dev/mapper/dev-test /media/wrong ext4 rw 0 0\n"),
					},
				},
			},
			wantErr: true,
		},
		"Fail to remove directory": {
			target:    "/media/vol",
			activator: testActivator{},
			mounter:   testMounter{},
			filesystem: testFileSystem{
				wantErr: true,
				MapFS: fstest.MapFS{
					"proc/mounts": &fstest.MapFile{
						Data: []byte("/dev/mapper/dev-test /media/vol ext4 rw 0 0\n"),
					},
				},
			},
			wantErr: true,
		},
		"Fail to unmount volume": {
			target:    "/media/vol",
			activator: testActivator{},
			mounter:   testMounter{wantErr: true},
			filesystem: testFileSystem{
				MapFS: fstest.MapFS{
					"proc/mounts": &fstest.MapFile{
						Data: []byte("/dev/mapper/dev-test /media/vol ext4 rw 0 0\n"),
					},
				},
			},
			wantErr: true,
		},
		"Fail to deactivate volume": {
			target:    "/media/vol",
			activator: testActivator{wantErr: true},
			mounter:   testMounter{},
			filesystem: testFileSystem{
				MapFS: fstest.MapFS{
					"proc/mounts": &fstest.MapFile{
						Data: []byte("/dev/mapper/dev-test /media/vol ext4 rw 0 0\n"),
					},
				},
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			ctx := testutils.ContextLoggerWithDebug(t)

			m := tpm.NewMount(
				tpmtestutils.WithActivator(&tc.activator),
				tpmtestutils.WithMounter(&tc.mounter),
				tpmtestutils.WithFileSystem(tc.filesystem),
			)

			err := m.UnmountVolume(ctx, tc.target)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.Equal(tc.activator.deactivated, tc.wantDectivated) // the volume is activated as expected
			is.Equal(tc.mounter.unmounted, tc.wantUnmounted)      // the volume is mounted as expected
		})
	}
}

type testActivator struct {
	activated   bool
	deactivated bool

	wantErr bool
}

func (m *testActivator) ActivateVolume(volumeName, device string, authRequestor secboot.AuthRequestor) error {
	if m.wantErr {
		return fmt.Errorf("test error")
	}

	m.activated = true
	return nil
}

func (m *testActivator) DeactivateVolume(volumeName string) error {
	if m.wantErr {
		return fmt.Errorf("test error")
	}

	m.deactivated = true
	return nil
}

type testFileSystem struct {
	fstest.MapFS

	readErr bool
	wantErr bool
}

func (fs testFileSystem) MkdirAll(path string) error {
	if fs.wantErr {
		return fmt.Errorf("test error")
	}
	return nil
}

func (fs testFileSystem) Open(name string) (fs.File, error) {
	if fs.readErr {
		return &errorFile{err: errors.New("simulated I/O error")}, nil
	}
	return fs.MapFS.Open(name)
}

func (fs testFileSystem) RemoveAll(path string) error {
	if fs.wantErr {
		return fmt.Errorf("test error")
	}
	return nil
}

type testMounter struct {
	mounted   bool
	unmounted bool

	wantErr bool
}

func (m *testMounter) Mount(path, target string) error {
	if m.wantErr {
		return fmt.Errorf("test error")
	}

	m.mounted = true
	return nil
}

func (m *testMounter) Unmount(target string) error {
	if m.wantErr {
		return fmt.Errorf("test error")
	}

	m.unmounted = true
	return nil
}

type errorFile struct{ err error }

func (e *errorFile) Read(p []byte) (int, error) { return 0, e.err }
func (e *errorFile) Close() error               { return nil }
func (e *errorFile) Stat() (fs.FileInfo, error) { return nil, errors.New("no stat") }
