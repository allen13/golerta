package main

import (
  "fmt"

  "github.com/dgrijalva/jwt-go"
  jwtmiddleware "github.com/iris-contrib/middleware/jwt"
  "github.com/kataras/iris"
  "github.com/allen13/golerta/auth"
  "github.com/docopt/docopt-go"
  "github.com/allen13/golerta/config"
)

const version = "Golerta 0.0.1"
const usage = `Golerta.

Usage:
  golerta --config=<config>
  golerta --help
  golerta --version

Options:
  --config=<config>            The golerta config [default: /etc/golerta/golerta.conf].
  --help                       Show this screen.
  --version                    Show version.
`

func main() {
  args, _ := docopt.Parse(usage, nil, true, version, false)
  configFile := args["--config"].(string)
  config := config.BuildConfig(configFile)
  fmt.Println(config)


  golerta := iris.New()

  authorizationMiddleware := buildAuthorizationMiddleware(config.Golerta.SigningKey)

  auth.RegisterAuthProvider(config, golerta)
  golerta.Get("/ping", pingHandler)
  golerta.Get("/secured/ping", authorizationMiddleware.Serve, securedPingHandler)

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

func buildAuthorizationMiddleware(signingKey string)(*jwtmiddleware.Middleware){
  return jwtmiddleware.New(jwtmiddleware.Config{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      return []byte(signingKey), nil
    },
    SigningMethod: jwt.SigningMethodHS256,
  })
}
