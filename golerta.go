package main

import (
  "github.com/dgrijalva/jwt-go"
  jwtmiddleware "github.com/iris-contrib/middleware/jwt"
  "github.com/kataras/iris"
  "github.com/allen13/golerta/app/auth"
  "github.com/docopt/docopt-go"
  "github.com/allen13/golerta/app/config"
  "github.com/allen13/golerta/app"
  "github.com/prometheus/common/log"
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

  golerta := iris.New()

  if config.Golerta.SigningKey == "" {
    log.Fatal("Shutting down, signing key must be provided.")
  }

  authorizationMiddleware := buildAuthorizationMiddleware(config.Golerta.SigningKey)

  auth.RegisterAuthProvider(config, golerta)
  app.RegisterApp(config,golerta,authorizationMiddleware)

  golerta.Listen(":8080")
}

func buildAuthorizationMiddleware(signingKey string)(*jwtmiddleware.Middleware){
  return jwtmiddleware.New(jwtmiddleware.Config{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      return []byte(signingKey), nil
    },
    SigningMethod: jwt.SigningMethodHS256,
  })
}
