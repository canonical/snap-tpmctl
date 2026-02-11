package snapd_test

import (
	"testing"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/matryer/is"
)

// Mock for snapd client backend (http)

func TestNew(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	s := snapd.New()
	is.True(s != nil) // New returned an object
}

