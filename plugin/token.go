package plugin

// Get token, token will be passed from client to server.
type TokenGetter interface {
	// You should cache the token, avoid generate new one every time.
	Get() ([]byte, error)
}

// Check token, token will be checked on server.
type TokenChecker interface {
	// You should not do a lot of time-consuming operations.
	Check([]byte) error
}
