package snapd_test

import (
	"testing"

	"snap-tpmctl/internal/snapd"
)

func TestFoo(t *testing.T) {
	snapd.WithInteraction(true)
}
