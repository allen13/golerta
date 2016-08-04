package auth

import (
	"github.com/kataras/iris"
)

//Interface for authentication providers
type AuthProvider interface {
	LoginHandler(ctx *iris.Context)
	SetSigningKey(key string)
	Connect() error
	Close()
}
