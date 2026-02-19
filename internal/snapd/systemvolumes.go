package snapd

import (
	"context"
	"encoding/json"
	"net/http"

	snapdClient "github.com/snapcore/snapd/client"
)

type SystemVolumesResult = snapdClient.SystemVolumesResult
type SystemVolumesStructureInfo = snapdClient.SystemVolumesStructureInfo
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
