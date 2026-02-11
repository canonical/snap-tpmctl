package tpm

import (
	"context"
	"fmt"

	"github.com/canonical/snap-tpmctl/internal/snapd"
)

type fdeStatusClient interface {
	FdeStatus(ctx context.Context) (*snapd.FdeStatusResult, error)
}

// FdeStatus retrieves the Full Disk Encryption status from snapd.
func FdeStatus(ctx context.Context, client fdeStatusClient) (string, error) {
	status, err := client.FdeStatus(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve the FDE status: %w", err)
	}

	return status.Status, nil
}
