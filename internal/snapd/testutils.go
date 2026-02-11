package snapd

// withSocketPath configures the snapd socket path for the client.
func withSocketPath(p string) Option {
	return func(o *options) {
		o.socket = p
	}
}
