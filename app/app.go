package app

import (
	"log"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"

	"github.com/allen13/golerta/app/auth"
	"github.com/allen13/golerta/app/auth/middleware"
	"github.com/allen13/golerta/app/config"
	"github.com/allen13/golerta/app/controllers"
	"github.com/allen13/golerta/app/services"
)

func BuildApp(config config.GolertaConfig) (e *echo.Echo) {
	config.Notifiers.Init()
	config.FlapDetection.Init()

	err := config.Rethinkdb.Init()
	if err != nil {
		log.Fatal(err)
	}
	db := config.Rethinkdb

	continuousQueryService := &services.ContinuousQueryService{
		DB:            db,
		QueryInterval: config.App.ContinuousQueryInterval.Duration,
		Notifiers:     config.Notifiers,
		FlapDetection: &config.FlapDetection,
	}
	go continuousQueryService.Start()

	e = echo.New()

	// enable CORS if UI is detached from the golerta process
	e.Use(echoMiddleware.CORS())

	authProvider := BuildAuthProvider(config)

	shouldSkipAuth := func(_ echo.Context) bool {
		c := config.App.AuthProvider == "noop"
		return c
	}

	authMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:     shouldSkipAuth,
		SigningKey:  []byte(config.App.SigningKey),
		TokenLookup: "header:Authorization,query:api-key",
	})

	authController := controllers.AuthController{
		Echo:         e,
		AuthProvider: authProvider,
	}
	authController.Init()

	alertsService := services.AlertService{
		DB:            &db,
		FlapDetection: &config.FlapDetection,
	}
	alertsController := controllers.AlertsController{
		Echo:             e,
		AlertService:     alertsService,
		AuthMiddleware:   authMiddleware,
		LogAlertRequests: config.App.LogAlertRequests,
	}
	alertsController.Init()

	e.Static("/static", "static")
	e.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(301, "/static/index.html")
	})

	return
}

func BuildAuthProvider(config config.GolertaConfig) (authProvider auth.AuthProvider) {
	switch config.App.AuthProvider {
	case "ldap":
		authProvider = &config.Ldap
	case "oauth":
		authProvider = &config.OAuth
	case "noop":
		authProvider = &config.Noop
	}

	err := authProvider.Connect()
	defer authProvider.Close()

	if err != nil {
		log.Fatal(err)
	}
	if config.App.SigningKey == "" {
		log.Fatal("Shutting down, signing key must be provided.")
	}
	authProvider.SetSigningKey(config.App.SigningKey)
	return
}
