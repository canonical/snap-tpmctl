// Package testutils provides testing utilities and mock implementations.
package testutils

import (
	"context"
	"encoding/json"
	"errors"

	"snap-tpmctl/internal/snapd"
)

// MockConfig holds configuration for MockSnapdClient behavior.
type MockConfig struct {
	// Error flags (for exceptional cases like snapd being down)
	LoadAuthError          bool
	GenerateKeyError       bool
	EnumerateError         bool
	AddKeyError            bool
	CheckPassphraseError   bool
	CheckPINError          bool
	ReplacePassphraseError bool
	ReplacePINError        bool

	// Validation error flags for CheckPassphrase
	PassphraseLowEntropy   bool
	PassphraseInvalid      bool
	PassphraseUnsupported  bool
	PassphraseUnknownError bool
	PassphraseNotOK        bool

	// Validation error flags for CheckPIN
	PINLowEntropy  bool
	PINInvalid     bool
	PINUnsupported bool
	PINNotOK       bool

	// Replace operation flags
	ReplacePassphraseNotOK bool
	ReplacePINNotOK        bool
}

// MockSnapdClient is a mock implementation of the snapdClienter interface for testing.
type MockSnapdClient struct {
	config MockConfig

	// Return values
	generatedKey  *snapd.GenerateRecoveryKeyResult
	systemVolumes *snapd.SystemVolumesResult
	asyncResp     *snapd.AsyncResponse
}

// NewMockSnapdClient creates a new mock snapd client with the given configuration.
func NewMockSnapdClient(cfg MockConfig) *MockSnapdClient {
	return &MockSnapdClient{
		config: cfg,
		generatedKey: &snapd.GenerateRecoveryKeyResult{
			KeyID:       "test-key-id-12345",
			RecoveryKey: "12345-67890-12345-67890-12345-67890-12345-67890",
		},
		systemVolumes: &snapd.SystemVolumesResult{
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
		asyncResp: &snapd.AsyncResponse{
			ID:      "change-123",
			Status:  "Done",
			Ready:   true,
			Summary: "Add recovery key",
		},
	}
}

// LoadAuthFromHome simulates loading authentication from the user's home directory.
func (m MockSnapdClient) LoadAuthFromHome() error {
	if m.config.LoadAuthError {
		return errors.New("mocked error for LoadAuthFromHome: cannot load auth: auth.json not found")
	}
	return nil
}

// GenerateRecoveryKey simulates generating a new recovery key.
func (m MockSnapdClient) GenerateRecoveryKey(ctx context.Context) (*snapd.GenerateRecoveryKeyResult, error) {
	if m.config.GenerateKeyError {
		return nil, errors.New("mocked error for GenerateRecoveryKey: cannot generate recovery key: snapd error")
	}
	return m.generatedKey, nil
}

// EnumerateKeySlots simulates enumerating system volume key slots.
func (m MockSnapdClient) EnumerateKeySlots(ctx context.Context) (*snapd.SystemVolumesResult, error) {
	if m.config.EnumerateError {
		return nil, errors.New("mocked error for EnumerateKeySlots: cannot enumerate key slots: snapd error")
	}
	return m.systemVolumes, nil
}

// AddRecoveryKey simulates adding a recovery key to specified slots.
func (m MockSnapdClient) AddRecoveryKey(ctx context.Context, keyID string, slots []snapd.KeySlot) (*snapd.AsyncResponse, error) {
	if m.config.AddKeyError {
		return nil, errors.New("mocked error for AddRecoveryKey: cannot add recovery key: permission denied")
	}
	return m.asyncResp, nil
}

// Close closes the mock client connection.
func (m MockSnapdClient) Close() error {
	return nil
}

// CheckPassphrase simulates checking if a passphrase is valid.
func (m MockSnapdClient) CheckPassphrase(ctx context.Context, passphrase string) (*snapd.Response, error) {
	if m.config.CheckPassphraseError {
		return nil, errors.New("mocked error for CheckPassphrase: cannot check passphrase: snapd error")
	}

	if m.config.PassphraseLowEntropy {
		return nil, &snapd.Error{
			Kind:    "invalid-passphrase",
			Message: "Mocked error for CheckPassphrase: passphrase is invalid",
			Value:   mustMarshalJSONForMock(map[string]any{"reasons": []string{"low-entropy"}, "entropy-bits": 24, "min-entropy-bits": 60, "optimal-entropy-bits": 80}),
		}
	}

	if m.config.PassphraseInvalid {
		return nil, &snapd.Error{
			Kind:    "invalid-passphrase",
			Message: "Mocked error for CheckPassphrase: passphrase contains invalid characters",
		}
	}

	if m.config.PassphraseUnsupported {
		return nil, &snapd.Error{
			Kind:    "unsupported",
			Message: "Mocked error for CheckPassphrase: passphrase validation is not available",
		}
	}

	if m.config.PassphraseUnknownError {
		return nil, &snapd.Error{
			Kind:    "unknown-error",
			Message: "Mocked error for CheckPassphrase: something went wrong",
		}
	}

	if m.config.PassphraseNotOK {
		return &snapd.Response{Status: "Bad Request", StatusCode: 400}, nil
	}

	return &snapd.Response{Status: "OK", StatusCode: 200}, nil
}

// CheckPIN simulates checking if a PIN is valid.
func (m MockSnapdClient) CheckPIN(ctx context.Context, pin string) (*snapd.Response, error) {
	if m.config.CheckPINError {
		return nil, errors.New("mocked error for CheckPIN: cannot check PIN: snapd error")
	}

	if m.config.PINLowEntropy {
		return nil, &snapd.Error{
			Kind:    "invalid-pin",
			Message: "Mocked error for CheckPIN: PIN is invalid",
			Value:   mustMarshalJSONForMock(map[string]any{"reasons": []string{"low-entropy"}, "entropy-bits": 13, "min-entropy-bits": 20, "optimal-entropy-bits": 30}),
		}
	}

	if m.config.PINInvalid {
		return nil, &snapd.Error{
			Kind:    "invalid-pin",
			Message: "Mocked error for CheckPIN: PIN format is invalid",
		}
	}

	if m.config.PINUnsupported {
		return nil, &snapd.Error{
			Kind:    "unsupported",
			Message: "Mocked error for CheckPIN: PIN validation is not available",
		}
	}

	if m.config.PINNotOK {
		return &snapd.Response{Status: "Bad Request", StatusCode: 400}, nil
	}

	return &snapd.Response{Status: "OK", StatusCode: 200}, nil
}

// ReplacePassphrase simulates replacing a passphrase.
func (m MockSnapdClient) ReplacePassphrase(ctx context.Context, oldPassphrase string, newPassphrase string, keySlots []snapd.KeySlot) (*snapd.AsyncResponse, error) {
	if m.config.ReplacePassphraseError {
		return nil, errors.New("mocked error for ReplacePassphrase: cannot replace passphrase: permission denied")
	}
	if m.config.ReplacePassphraseNotOK {
		return &snapd.AsyncResponse{
			ID:     "change-123",
			Status: "Error",
			Ready:  false,
		}, nil
	}
	return m.asyncResp, nil
}

// ReplacePIN simulates replacing a PIN.
func (m MockSnapdClient) ReplacePIN(ctx context.Context, oldPin string, newPin string, keySlots []snapd.KeySlot) (*snapd.AsyncResponse, error) {
	if m.config.ReplacePINError {
		return nil, errors.New("mocked error for ReplacePIN: cannot replace PIN: permission denied")
	}
	if m.config.ReplacePINNotOK {
		return &snapd.AsyncResponse{
			ID:     "change-123",
			Status: "Error",
			Ready:  false,
		}, nil
	}
	return m.asyncResp, nil
}

// mustMarshalJSONForMock marshals a value to JSON for use in mock responses.
func mustMarshalJSONForMock(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic("failed to marshal JSON in mock: " + err.Error())
	}
	return data
}
