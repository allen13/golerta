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

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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
	standardResponse(ctx, alertsResponse, err)
}

func (ac *AlertsController) getAlert(ctx *iris.Context) {
	alertResponse, err := ac.AlertService.GetAlert(ctx.Param("alert"))
	standardResponse(ctx, alertResponse, err)
}

func (ac *AlertsController) getAlertsCount(ctx *iris.Context) {
	alertsCount, err := ac.AlertService.GetAlertsCount(ctx.QueryArgs())
	standardResponse(ctx, alertsCount, err)
}

func (ac *AlertsController) getAlertsServices(ctx *iris.Context) {
	groupedServices, err := ac.AlertService.GetGroupedServices(ctx.QueryArgs())
	standardResponse(ctx, groupedServices, err)
}

func (ac *AlertsController) getAlertsEnvironments(ctx *iris.Context) {
	groupedEnvironments, err := ac.AlertService.GetGroupedEnvironments(ctx.QueryArgs())
	standardResponse(ctx, groupedEnvironments, err)
}

func standardResponse(ctx *iris.Context, response interface{}, err error){
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, ErrorResponse{Status: "error", Message: err.Error()})
	}
	ctx.JSON(iris.StatusOK, response)
}
