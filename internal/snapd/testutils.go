package snapd

import "github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"

// withBaseURL configures the snapd server to connect to.
func withBaseURL(p string) Option {
	testsdetection.MustBeTesting()
	return func(o *options) {
		o.baseURL = p
	}
}
