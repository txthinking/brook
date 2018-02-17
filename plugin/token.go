package plugin

// Token can do additional certification
type Token interface {
	// Get token, token will be passed from client to server.
	// You should cache the token, avoid generate new one every time.
	Get() ([]byte, error)

	// Check token, token will be checked on server.
	// You should not do a lot of time-consuming operations.
	Check([]byte) error
}
