// Package tpm manages TPM/FDE features
package tpm

import (
	"context"
	"fmt"
	_ "unsafe" // Needed for go:linkname.

	"github.com/canonical/snap-tpmctl/internal/snapd"
)

type SnapTPM struct {
	snapdClient *snapd.Client
}

func New() SnapTPM {
	return SnapTPM{snapdClient: snapd.New()}
}

// FdeStatus retrieves the Full Disk Encryption status from snapd.
func (s SnapTPM) FdeStatus(ctx context.Context) (string, error) {
	status, err := s.snapdClient.FdeStatus(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve the FDE status: %v", err)
	}

	return status, nil
}
