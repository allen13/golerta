package noop

type NoopAuthProvider struct {
	signingKey string
}

func (na *NoopAuthProvider) SetSigningKey(_ string) {
	na.signingKey = ""
}

// Connect bypassed for no auth provider
func (na *NoopAuthProvider) Connect() error {
	return nil
}

// Close bypassed for no auth provider
func (na *NoopAuthProvider) Close() {
	return
}

// Authenticate returns true to noop the auth provider
func (na *NoopAuthProvider) Authenticate(_ string, _ string) (authenticated bool, token string, err error) {
	token = ""
	authenticated = true
	err = nil
	return
}
