package auth

import (
  "github.com/kataras/iris"
  "github.com/allen13/golerta/config"
)

//Interface for authentication providers
type AuthProvider interface {
  LoginHandler(ctx *iris.Context)
  SetSigningKey(key string)
  Connect()(error)
  Close()
}

func RegisterAuthProvider(config config.GolertaConfig, golerta *iris.Framework){
  var authProvider AuthProvider

  switch config.Golerta.AuthProvider{
  case "ldap":
   authProvider = &config.Ldap
  }

  authProvider.Connect()

  golerta.Post("/auth/login", authProvider.LoginHandler)
}