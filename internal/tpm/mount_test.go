package tpm_test

import (
	"errors"
	"io"
	"io/fs"
	"testing"
	"testing/fstest"
	"testing/iotest"

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

		filesystem testFileSystem
		volume     testVolume

		wantActivated bool
		wantMounted   bool

		wantErr bool
	}{
		"Success on mounting volume": {
			device:        "/dev/test",
			target:        "/media/vol",
			volume:        testVolume{},
			filesystem:    testFileSystem{},
			wantActivated: true,
			wantMounted:   true,
		},
		"Success on mounting already active volume": {
			device: "/dev/test",
			target: "/media/vol",
			volume: testVolume{},
			filesystem: testFileSystem{
				MapFS: fstest.MapFS{
					"dev/mapper/dev-test": &fstest.MapFile{},
				},
			},
			wantMounted: true,
		},

		"Fail to create directory": {
			device: "/dev/test",
			target: "/media/vol",
			volume: testVolume{},
			filesystem: testFileSystem{
				wantErr: true,
			},
			wantErr: true,
		},
		"Fail to activate volume": {
			device:     "/dev/test",
			target:     "/media/vol",
			volume:     testVolume{wantActivateErr: true},
			filesystem: testFileSystem{},
			wantErr:    true,
		},
		"Fail to mount volume": {
			device:        "/dev/test",
			target:        "/media/vol",
			volume:        testVolume{wantMountErr: true},
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
				tpmtestutils.WithVolume(&tc.volume),
				tpmtestutils.WithFileSystem(tc.filesystem),
			)

			err := m.MountVolume(ctx, tc.device, tc.target)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.Equal(tc.volume.activated, tc.wantActivated) // the volume is activated as expected
			is.Equal(tc.volume.mounted, tc.wantMounted)     // the volume is mounted as expected
		})
	}
}

func TestUnmountVolume(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		target string

		filesystem testFileSystem
		volume     testVolume

		wantDectivated bool
		wantUnmounted  bool

		wantErr bool
	}{
		"Success on deactivating volume": {
			target: "/media/vol",
			volume: testVolume{},
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

		"Fail to remove directory": {
			target: "/media/vol",
			volume: testVolume{},
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
			target: "/media/vol",
			volume: testVolume{wantMountErr: true},
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
			target: "/media/vol",
			volume: testVolume{wantActivateErr: true},
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
				tpmtestutils.WithVolume(&tc.volume),
				tpmtestutils.WithFileSystem(tc.filesystem),
			)

			err := m.UnmountVolume(ctx, tc.target)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.Equal(tc.volume.deactivated, tc.wantDectivated) // the volume is deactivated as expected
			is.Equal(tc.volume.unmounted, tc.wantUnmounted)    // the volume is unmounted as expected
		})
	}
}

func TestGetMapperFromMount(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		mapper string
		target string

		filesystem testFileSystem

		wantErr bool
	}{
		"Success on getting mapper": {
			mapper: "/dev/mapper/dev-test",
			target: "/media/vol",
			filesystem: testFileSystem{
				MapFS: fstest.MapFS{
					"proc/mounts": &fstest.MapFile{
						Data: []byte("/dev/mapper/dev-test /media/vol ext4 rw 0 0\n"),
					},
				},
			},
		},

		"Fail to get mapper /proc/mounts": {
			target:     "/media/vol",
			filesystem: testFileSystem{},
			wantErr:    true,
		},
		"Fail to read /proc/mounts": {
			target: "/media/vol",
			filesystem: testFileSystem{
				wantReadErr: true,
			},
			wantErr: true,
		},
		"Fail to find mapper from /proc/mounts": {
			target: "/media/vol",
			filesystem: testFileSystem{
				MapFS: fstest.MapFS{
					"proc/mounts": &fstest.MapFile{
						Data: []byte("/dev/mapper/dev-test /media/wrong ext4 rw 0 0\n"),
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

			m := tpm.NewMount(
				tpmtestutils.WithFileSystem(tc.filesystem),
			)

			mapper, err := tpm.GetMapperFromMount(m, tc.target)
			if tc.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			is.Equal(mapper, tc.mapper) // the device mapper is the expected one
		})
	}
}

type testFileSystem struct {
	fstest.MapFS

	wantReadErr bool
	wantErr     bool
}

func (fs testFileSystem) MkdirAll(path string) error {
	if fs.wantErr {
		return errors.New("test error")
	}
	return nil
}

func (fs testFileSystem) Open(name string) (fs.File, error) {
	if fs.wantReadErr {
		return &errorFile{iotest.ErrReader(errors.New("simulated I/O error"))}, nil
	}
	return fs.MapFS.Open(name)
}

func (fs testFileSystem) RemoveAll(path string) error {
	if fs.wantErr {
		return errors.New("test error")
	}
	return nil
}

type testVolume struct {
	activated   bool
	deactivated bool
	mounted     bool
	unmounted   bool

	wantActivateErr bool
	wantMountErr    bool
}

func (m *testVolume) Activate(volumeName, device string, authRequestor secboot.AuthRequestor) error {
	if m.wantActivateErr {
		return errors.New("test error")
	}

	m.activated = true
	return nil
}

func (m *testVolume) Deactivate(volumeName string) error {
	if m.wantActivateErr {
		return errors.New("test error")
	}

	m.deactivated = true
	return nil
}

func (m *testVolume) Mount(path, target string) error {
	if m.wantMountErr {
		return errors.New("test error")
	}

	m.mounted = true
	return nil
}

func (m *testVolume) Unmount(target string) error {
	if m.wantMountErr {
		return errors.New("test error")
	}

	m.unmounted = true
	return nil
}

type errorFile struct{ io.Reader }

func (e *errorFile) Close() error               { return nil }
func (e *errorFile) Stat() (fs.FileInfo, error) { return nil, errors.New("err stat") }
