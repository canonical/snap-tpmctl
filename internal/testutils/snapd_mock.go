package testutils

import (
	"context"
	"errors"

	"snap-tpmctl/internal/snapd"
)
// FIXME: move this into a submodule, testutils_snapd/
// MockSnapdClient is a mock implementation of the snapdClienter interface for testing.
type MockSnapdClient struct {
	// Error flags
	LoadAuthError    bool
	GenerateKeyError bool
	EnumerateError   bool
	AddKeyError      bool

	// Return values
	GeneratedKey  *snapd.GenerateRecoveryKeyResult
	SystemVolumes *snapd.SystemVolumesResult
	AsyncResp     *snapd.AsyncResponse
}

// NewMockSnapdClient creates a new mock snapd client with default success responses.
func NewMockSnapdClient() *MockSnapdClient {
	return &MockSnapdClient{
		GeneratedKey: &snapd.GenerateRecoveryKeyResult{
			KeyID:       "test-key-id-12345",
			RecoveryKey: "12345-67890-12345-67890-12345-67890-12345-67890",
		},
		SystemVolumes: &snapd.SystemVolumesResult{
			ByContainerRole: map[string]snapd.VolumeInfo{
				"system-data": {
					Name:       "ubuntu-data",
					VolumeName: "pc",
					Encrypted:  true,
					KeySlots: map[string]snapd.KeySlotInfo{
						"default": {
							Type:         "platform",
							AuthMode:     "passphrase",
							PlatformName: "tpm2",
							Roles:        []string{"run+recover"},
						},
						"default-fallback": {
							Type:         "platform",
							AuthMode:     "passphrase",
							PlatformName: "tpm2",
							Roles:        []string{"recover"},
						},
						"default-recovery": {
							Type: "recovery",
						},
						"additional-recovery": {
							Type: "recovery",
						},
					},
				},
				"system-save": {
					Name:       "ubuntu-save",
					VolumeName: "pc",
					Encrypted:  true,
					KeySlots: map[string]snapd.KeySlotInfo{
						"default": {
							Type:         "platform",
							AuthMode:     "none",
							PlatformName: "plainkey",
						},
						"default-fallback": {
							Type:         "platform",
							AuthMode:     "passphrase",
							PlatformName: "tpm2",
							Roles:        []string{"recover"},
						},
						"default-recovery": {
							Type: "recovery",
						},
						"additional-recovery": {
							Type: "recovery",
						},
					},
				},
			},
		},
		AsyncResp: &snapd.AsyncResponse{
			ID:      "change-123",
			Status:  "Done",
			Ready:   true,
			Summary: "Add recovery key",
		},
	}
}

func (m MockSnapdClient) LoadAuthFromHome() error {
	if m.LoadAuthError {
		return errors.New("cannot load auth: auth.json not found")
	}
	return nil
}

func (m MockSnapdClient) GenerateRecoveryKey(ctx context.Context) (*snapd.GenerateRecoveryKeyResult, error) {
	if m.GenerateKeyError {
		return nil, errors.New("cannot generate recovery key: snapd error")
	}
	return m.GeneratedKey, nil
}

func (m MockSnapdClient) EnumerateKeySlots(ctx context.Context) (*snapd.SystemVolumesResult, error) {
	if m.EnumerateError {
		return nil, errors.New("cannot enumerate key slots: snapd error")
	}
	return m.SystemVolumes, nil
}

func (m MockSnapdClient) AddRecoveryKey(ctx context.Context, keyID string, slots []snapd.KeySlot) (*snapd.AsyncResponse, error) {
	if m.AddKeyError {
		return nil, errors.New("cannot add recovery key: permission denied")
	}
	return m.AsyncResp, nil
}

func (m MockSnapdClient) Close() error {
	return nil
}
