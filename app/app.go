package app

import (
  "github.com/kataras/iris"
  "github.com/allen13/golerta/app/config"
  jwtmiddleware "github.com/iris-contrib/middleware/jwt"
)

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

func RegisterApp(config config.GolertaConfig, golerta *iris.Framework,authorizationMiddleware *jwtmiddleware.Middleware){
  golerta.Get("/ping", pingHandler)
  golerta.Get("/secured/ping", authorizationMiddleware.Serve, securedPingHandler)
}
