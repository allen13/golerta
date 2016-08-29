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
	"github.com/kataras/iris/utils"
)

func BuildApp(config config.GolertaConfig) (http *iris.Framework) {
	config.Notifiers.Init()
	config.FlapDetection.Init()

	err := config.Rethinkdb.Init()
	if err != nil {
		log.Fatal(err)
	}
	db := config.Rethinkdb

	continuousQueryService :=	&services.ContinuousQueryService{
		DB: db,
		QueryInterval: config.Golerta.ContinuousQueryInterval.Duration,
		Notifiers: config.Notifiers,
		FlapDetection: &config.FlapDetection,
	}
	go continuousQueryService.Start()

	http = iris.New()

	BuildAuthProvider(config, http)
	authMiddleware := BuildAuthorizationMiddleware(config.Golerta.SigningKey)

	alertsService := services.AlertService{
		DB: &db,
		FlapDetection: &config.FlapDetection,
	}
	alertsController := controllers.AlertsController{
		HTTP:           http,
		AlertService:   alertsService,
		AuthMiddleware: authMiddleware,
	}
	alertsController.Init()

	StaticWeb(http, "/static", "./static", 1)
	http.Get("/", func(ctx *iris.Context) {
		ctx.Redirect("/static/index.html", 301)
	})

	return
}

func StaticWeb(http *iris.Framework, reqPath string, systemPath string, stripSlashes int) {
	if reqPath[len(reqPath)-1] != byte('/') { // if / then /*filepath, if /something then /something/*filepath
		reqPath += "/"
	}

	hasIndex := utils.Exists(systemPath + utils.PathSeparator + "index.html")
	serveHandler := http.StaticHandler(systemPath, stripSlashes, false, !hasIndex, nil) // if not index.html exists then generate index.html which shows the list of files
	indexHandler := func(ctx *iris.Context) {
		if len(ctx.Param("filepath")) < 2 && hasIndex {
			ctx.Request.SetRequestURI(reqPath + "index.html")
		}
		ctx.Next()

	}
	http.Head(reqPath+"*filepath", indexHandler, serveHandler)
	http.Get(reqPath+"*filepath", indexHandler, serveHandler)
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
