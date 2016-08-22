package app

import (
	"github.com/allen13/golerta/app/auth"
	"github.com/allen13/golerta/app/controllers"
	"github.com/allen13/golerta/app/services"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"log"
	"github.com/allen13/golerta/app/config"
)

func BuildApp(config config.GolertaConfig) (http *iris.Framework) {
	err := config.Rethinkdb.Init()
	if err != nil {
		log.Fatal(err)
	}
	db := config.Rethinkdb

	continuousQueryService :=	&services.ContinuousQueryService{DB: db, QueryInterval: config.Golerta.ContinuousQueryInterval.Duration}
	go continuousQueryService.Start()

	http = iris.New()

	BuildAuthProvider(config, http)
	authMiddleware := BuildAuthorizationMiddleware(config.Golerta.SigningKey)

	alertsService := services.AlertService{DB: db}
	alertsController := controllers.AlertsController{
		HTTP:           http,
		AlertService:   alertsService,
		AuthMiddleware: authMiddleware,
	}
	alertsController.Init()

	http.StaticWeb("/static", "./static", 1)
	http.Get("/", func(ctx *iris.Context) {
		ctx.Redirect("/static/index.html", 301)
	})

	return
}

func BuildAuthProvider(config config.GolertaConfig, http *iris.Framework) {
	var authProvider auth.AuthProvider
	switch config.Golerta.AuthProvider {
	case "ldap":
		authProvider = &config.Ldap
	}

	err := authProvider.Connect()
	defer authProvider.Close()

	if err != nil {
		log.Fatal(err)
	}
	if config.Golerta.SigningKey == "" {
		log.Fatal("Shutting down, signing key must be provided.")
	}
	authProvider.SetSigningKey(config.Golerta.SigningKey)

	http.Post("/auth/login", authProvider.LoginHandler)
}

func BuildAuthorizationMiddleware(signingKey string) *jwtmiddleware.Middleware {
	jwtCofig := jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader, jwtmiddleware.FromParameter("api-key")),
	}
	return jwtmiddleware.New(jwtCofig)
}
