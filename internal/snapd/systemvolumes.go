package snapd

import (
	"context"
	"encoding/json"
	"net/http"

	snapdClient "github.com/snapcore/snapd/client"
)

// SystemVolumesResult is the response for /v2/system-volumes and includes the
// list of volumes plus their structures and keyslots.
type SystemVolumesResult = snapdClient.SystemVolumesResult

// SystemVolumesStructureInfo describes a single structure within a system
// volume (name, device, fs, and size info).
type SystemVolumesStructureInfo = snapdClient.SystemVolumesStructureInfo

// KeySlotInfo contains auth mode/type data used to distinguish recovery key,
// passphrase, or PIN-backed slots.
type KeySlotInfo = snapdClient.KeyslotInfo

func IsRecoveryKey(slot KeySlotInfo) bool {
	return slot.Type == snapdClient.KeyslotTypeRecovery
}

// IsPassphrase returns true if the keyslot uses passphrase authentication.
func IsPassphrase(slot KeySlotInfo) bool {
	return slot.AuthMode == AuthModePassphrase
}

// IsPIN returns true if the keyslot uses pin authentication.
func IsPIN(slot KeySlotInfo) bool {
	return slot.AuthMode == AuthModePIN
}

// ListVolumeInfo gets information about system volumes.
func (c *Client) ListVolumeInfo(ctx context.Context) (result SystemVolumesResult, err error) {
	resp, err := c.doSyncRequest(ctx, http.MethodGet, "/v2/system-volumes", nil, nil, nil)
	if err != nil {
		return result, err
	}

	var volumes SystemVolumesResult
	if err := json.Unmarshal(resp.Result, &volumes); err != nil {
		return result, err
	}

	return volumes, nil
}
