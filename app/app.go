package app

import (
  "log"
  "github.com/BurntSushi/toml"
  "github.com/allen13/golerta/app/services"
  "github.com/kataras/iris"
  "github.com/allen13/golerta/app/auth"
  "github.com/allen13/golerta/app/controllers"
  "github.com/allen13/golerta/app/db/rethinkdb"
  "github.com/allen13/golerta/app/auth/ldap"
  "github.com/dgrijalva/jwt-go"
  jwtmiddleware "github.com/iris-contrib/middleware/jwt"
)

type GolertaConfig struct {
  Golerta golerta
  Ldap ldap.LDAPAuthProvider
  Rethinkdb rethinkdb.RethinkDB
}

type golerta struct {
  SigningKey string `toml:"signing_key"`
  AuthProvider string `toml:"auth_provider"`
}

func BuildApp(configFile string)(http *iris.Framework){
  config := BuildConfig(configFile)

  err := config.Rethinkdb.Init()
  if err != nil{
    log.Fatal(err)
  }
  db := &config.Rethinkdb

  http = iris.New()

  BuildAuthProvider(config, http)
  authMiddleware := BuildAuthorizationMiddleware(config.Golerta.SigningKey)

  alertsService := services.AlertService{DB: db}
  alertsController := controllers.AlertsController{
    HTTP: http,
    AlertService: alertsService,
    AuthMiddleware: authMiddleware,
  }
  alertsController.Init()

  http.StaticWeb("/static", "./static", 1)
  http.Get("/",func(ctx *iris.Context) {
    ctx.Redirect("/static/index.html", 301)
  })

  return
}

func BuildConfig(configFile string)(config GolertaConfig){
  _, err := toml.DecodeFile(configFile, &config)

  if err != nil {
    log.Fatal("config file error: " + err.Error())
  }

  setDefaultConfigs(&config)
  return
}

func setDefaultConfigs(config* GolertaConfig){
  if config.Golerta.AuthProvider == ""{
    config.Golerta.AuthProvider = "ldap"
  }
}

func BuildAuthProvider(config GolertaConfig, http *iris.Framework){
  var authProvider auth.AuthProvider
  switch config.Golerta.AuthProvider{
  case "ldap":
    authProvider = &config.Ldap
  }

  err := authProvider.Connect()
  if err != nil {
    log.Fatal(err)
  }
  if config.Golerta.SigningKey == "" {
    log.Fatal("Shutting down, signing key must be provided.")
  }
  authProvider.SetSigningKey(config.Golerta.SigningKey)

  http.Post("/auth/login", authProvider.LoginHandler)
}

func BuildAuthorizationMiddleware(signingKey string)(*jwtmiddleware.Middleware){
  return jwtmiddleware.New(jwtmiddleware.Config{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      return []byte(signingKey), nil
    },
    SigningMethod: jwt.SigningMethodHS256,
  })
}
