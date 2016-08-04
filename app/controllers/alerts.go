package controllers

import (
	"github.com/allen13/golerta/app/services"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
)

type AlertsController struct {
	HTTP           *iris.Framework
	AlertService   services.AlertService
	AuthMiddleware *jwtmiddleware.Middleware
}

type error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (ac *AlertsController) Init() {
	ac.HTTP.Get("/alerts", ac.getAlerts)
	ac.HTTP.Get("/alert/:alert", ac.getAlert)
}

func (ac *AlertsController) getAlerts(ctx *iris.Context) {
	alertsResponse, err := ac.AlertService.GetAlerts(ctx.QueryArgs())
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, error{Status: "error", Message: err.Error()})
	}
	ctx.JSON(iris.StatusOK, alertsResponse)
}

func (ac *AlertsController) getAlert(ctx *iris.Context) {
	alertResponse, err := ac.AlertService.GetAlert(ctx.Param("alert"))
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, error{Status: "error", Message: err.Error()})
	}
	ctx.JSON(iris.StatusOK, alertResponse)
}
