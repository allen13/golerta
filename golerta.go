package main

import (
  "github.com/dgrijalva/jwt-go"
  jwtmiddleware "github.com/iris-contrib/middleware/jwt"
  "github.com/kataras/iris"
  "github.com/allen13/golerta/auth"
)

func main() {
  golerta := iris.New()
  signingKey := "AllYourBase"
  myJwtMiddleware := jwtmiddleware.New(jwtmiddleware.Config{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      return []byte(signingKey), nil
    },
    SigningMethod: jwt.SigningMethodHS256,
  })

  auth.RegisterAuthProvider("ldap", golerta)
  golerta.Get("/ping", pingHandler)
  golerta.Get("/secured/ping", myJwtMiddleware.Serve, securedPingHandler)

  golerta.Listen(":8080")
}

type response struct {
  Text string `json:"text"`
}

func pingHandler(ctx *iris.Context) {
  response := response{"All good. You don't need to be authenticated to call this"}
  ctx.JSON(iris.StatusOK, response)
}

func securedPingHandler(ctx *iris.Context) {
  response := response{"All good. You only get this message if you're authenticated"}
  // get the *jwt.Token which contains user information using:
  // user:= myJwtMiddleware.Get(ctx) or context.Get("jwt").(*jwt.Token)
  ctx.JSON(iris.StatusOK, response)
}
