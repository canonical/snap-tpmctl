package snapd_test

import (
	"snap-tpmctl/internal/snapd"
	"testing"
)

func TestFoo(t *testing.T) {
	snapd.WithInteraction(true)
}
