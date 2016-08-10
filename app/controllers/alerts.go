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
	ac.HTTP.Post("/alert", ac.AuthMiddleware.Serve, ac.createAlert)
	ac.HTTP.Get("/alert/:alert", ac.AuthMiddleware.Serve, ac.getAlert)
	ac.HTTP.Get("/alerts", ac.AuthMiddleware.Serve, ac.getAlerts)
	ac.HTTP.Post("/alert/:alert/status", ac.AuthMiddleware.Serve, ac.updateAlertStatus)
	ac.HTTP.Delete("/alert/:alert", ac.AuthMiddleware.Serve, ac.deleteAlert)
	ac.HTTP.Get("/alerts/count", ac.AuthMiddleware.Serve, ac.getAlertsCount)
	ac.HTTP.Get("/alerts/services", ac.AuthMiddleware.Serve, ac.getAlertsServices)
	ac.HTTP.Get("/alerts/environments", ac.AuthMiddleware.Serve, ac.getAlertsEnvironments)

}

func (ac *AlertsController) createAlert(ctx *iris.Context) {

	var incomingAlert models.Alert
	err := ctx.ReadJSON(&incomingAlert)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, models.ErrorResponse{Status: "error", Message: err.Error()})
	}

	alertsResponse, err := ac.AlertService.ProcessAlert(incomingAlert)
	models.StandardResponse(ctx, alertsResponse, err)
}

func (ac *AlertsController) getAlerts(ctx *iris.Context) {
	alertsResponse, err := ac.AlertService.GetAlerts(ctx.QueryArgs())
	models.StandardResponse(ctx, alertsResponse, err)
}

func (ac *AlertsController) getAlert(ctx *iris.Context) {
	alertResponse, err := ac.AlertService.GetAlert(ctx.Param("alert"))
	models.StandardResponse(ctx, alertResponse, err)
}

func (ac *AlertsController) deleteAlert(ctx *iris.Context) {
	err := ac.AlertService.DeleteAlert(ctx.Param("alert"))
	models.StandardResponse(ctx, struct {Status string `json:"status"`}{Status: "ok"}, err)
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

func (ac *AlertsController) updateAlertStatus(ctx *iris.Context) {
	var alertStatusUpdateRequest models.AlertStatusUpdateRequest
	err := ctx.ReadJSON(&alertStatusUpdateRequest)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, models.ErrorResponse{Status: "error", Message: err.Error()})
	}

	err = ac.AlertService.UpdateAlertStatus(ctx.Param("alert"), alertStatusUpdateRequest.Status, alertStatusUpdateRequest.Text)
	models.StandardResponse(ctx, iris.Map{"status":"ok"}, err)
}