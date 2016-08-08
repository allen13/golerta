package controllers

import (
	"github.com/allen13/golerta/app/services"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/allen13/golerta/app/models"
)

type AlertsController struct {
	HTTP           *iris.Framework
	AlertService   services.AlertService
	AuthMiddleware *jwtmiddleware.Middleware
}

func (ac *AlertsController) Init() {
	ac.HTTP.Get("/alert/:alert", ac.getAlert)
	ac.HTTP.Get("/alerts", ac.getAlerts)
	ac.HTTP.Get("/alerts/count", ac.getAlertsCount)
	ac.HTTP.Get("/alerts/services", ac.getAlertsServices)
	ac.HTTP.Get("/alerts/environments", ac.getAlertsEnvironments)
}

func (ac *AlertsController) getAlerts(ctx *iris.Context) {
	alertsResponse, err := ac.AlertService.GetAlerts(ctx.QueryArgs())
	models.StandardResponse(ctx, alertsResponse, err)
}

func (ac *AlertsController) getAlert(ctx *iris.Context) {
	alertResponse, err := ac.AlertService.GetAlert(ctx.Param("alert"))
	models.StandardResponse(ctx, alertResponse, err)
}

func (ac *AlertsController) getAlertsCount(ctx *iris.Context) {
	alertsCount, err := ac.AlertService.GetAlertsCount(ctx.QueryArgs())
	models.StandardResponse(ctx, alertsCount, err)
}

func (ac *AlertsController) getAlertsServices(ctx *iris.Context) {
	groupedServices, err := ac.AlertService.GetGroupedServices(ctx.QueryArgs())
	models.StandardResponse(ctx, groupedServices, err)
}

func (ac *AlertsController) getAlertsEnvironments(ctx *iris.Context) {
	groupedEnvironments, err := ac.AlertService.GetGroupedEnvironments(ctx.QueryArgs())
	models.StandardResponse(ctx, groupedEnvironments, err)
}
