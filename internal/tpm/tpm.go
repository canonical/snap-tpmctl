// Package tpm manages TPM/FDE features
package tpm

import (
	"context"
	"fmt"
	_ "unsafe" // Needed for go:linkname.

	"github.com/canonical/snap-tpmctl/internal/snapd"
)

// SnapTPM provides methods to interact with TPM/FDE features via snapd.
type SnapTPM struct {
	snapdClient *snapd.Client
}

type options struct {
	snapdClient *snapd.Client
}

// Option is a functional option for configuring the SnapTPM.
type Option func(*options)

// New creates a new SnapTPM instance with the provided options.
func New(args ...Option) SnapTPM {
	o := options{
		snapdClient: snapd.New(),
	}
	for _, f := range args {
		f(&o)
	}

	return SnapTPM{snapdClient: o.snapdClient}
}

// FdeStatus retrieves the Full Disk Encryption status from snapd.
func (s SnapTPM) FdeStatus(ctx context.Context) (string, error) {
	status, err := s.snapdClient.FdeStatus(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve the FDE status: %v", err)
	}

	return status, nil
}
