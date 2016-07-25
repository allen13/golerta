package auth

import (
  "github.com/kataras/iris"
)

//Interface for authentication providers
type AuthProvider interface {
  LoginHandler(ctx *iris.Context)
  SetSigningKey(key string)
  Connect()(error)
  Close()
}

//JSON struct that holds generated authorization token
type AuthToken struct {
  Token string `json:"token"`
}

//JSON struct for login errors
type LoginError struct {
  Status string  `json:"status"`
  Message string `json:"message"`
}

func RegisterAuthProvider(provider string, golerta *iris.Framework){
  var authProvider AuthProvider
  switch provider{
  case "ldap":
    authProvider = &LDAPAuthProvider{
      Host: "ldap.forumsys.com",
      Port: 389,
      BindDN: "cn=read-only-admin,dc=example,dc=com",
      BindPassword: "password",
      UserFilter:   "(uid=%s)",
      Base: "dc=example,dc=com",
    }
  }
  authProvider.Connect()
  defer authProvider.Close()

  golerta.Post("/auth/login", authProvider.LoginHandler)
}