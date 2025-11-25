package snapd

// WithSocketPath sets a custom socket path.
func WithSocketPath(path string) ClientOption {
	return func(c *Client) {
		c.socketPath = path
	}
}

// WithUserAgent sets a custom user agent.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithInteraction enables interactive authentication (Polkit dialogs).
func WithInteraction(allow bool) ClientOption {
	return func(c *Client) {
		c.allowInteraction = allow
	}
}
