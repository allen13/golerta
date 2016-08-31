package auth

//Interface for authentication providers
type AuthProvider interface {
	Authenticate(username, password string) (authenticated bool, token string, err error)
	SetSigningKey(key string)
	Connect() error
	Close()
}
