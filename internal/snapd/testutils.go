package snapd

// withBaseURL configures the snapd server to connect to.
func withBaseURL(p string) Option {
	return func(o *options) {
		o.baseURL = p
	}
}
