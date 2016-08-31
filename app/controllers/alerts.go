package controllers

import (
	"github.com/allen13/golerta/app/models"
	"github.com/allen13/golerta/app/services"
	"github.com/labstack/echo"
	"net/http"
)

type AlertsController struct {
	Echo           *echo.Echo
	AlertService   services.AlertService
	AuthMiddleware echo.MiddlewareFunc
}

func (ac *AlertsController) Init() {
	//ac.HTTP.Post("/alert", ac.AuthMiddleware.Serve, ac.createAlert)
	//ac.HTTP.Get("/alert/:alert", ac.AuthMiddleware.Serve, ac.getAlert)
	//ac.HTTP.Get("/alerts", ac.AuthMiddleware.Serve, ac.getAlerts)
	//ac.HTTP.Post("/alert/:alert/status", ac.AuthMiddleware.Serve, ac.updateAlertStatus)
	//ac.HTTP.Delete("/alert/:alert", ac.AuthMiddleware.Serve, ac.deleteAlert)
	//ac.HTTP.Get("/alerts/count", ac.AuthMiddleware.Serve, ac.getAlertsCount)
	//ac.HTTP.Get("/alerts/services", ac.AuthMiddleware.Serve, ac.getAlertsServices)
	//ac.Echo.Get("/alerts/environments", ac.AuthMiddleware.Serve, ac.getAlertsEnvironments)

	ac.Echo.Post("/alert", ac.createAlert, ac.AuthMiddleware)
	ac.Echo.Get("/alert/:alert", ac.getAlert, ac.AuthMiddleware)
	ac.Echo.Get("/alerts", ac.getAlerts, ac.AuthMiddleware)
	ac.Echo.Post("/alert/:alert/status", ac.updateAlertStatus, ac.AuthMiddleware)
	ac.Echo.Delete("/alert/:alert", ac.deleteAlert, ac.AuthMiddleware)
	ac.Echo.Get("/alerts/count", ac.getAlertsCount, ac.AuthMiddleware)
	ac.Echo.Get("/alerts/services", ac.getAlertsServices, ac.AuthMiddleware)
	ac.Echo.Get("/alerts/environments", ac.getAlertsEnvironments, ac.AuthMiddleware)

}

func (ac *AlertsController) createAlert(ctx echo.Context) error {
	ctx.Request().Body()
	var incomingAlert models.Alert
	err := ctx.Bind(&incomingAlert)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
	}

	alertsResponse, err := ac.AlertService.ProcessAlert(incomingAlert)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}

	return ctx.JSON(http.StatusCreated, alertsResponse)
	return nil
}

func (ac *AlertsController) getAlerts(ctx echo.Context) error {
	ctx.QueryParams()
	alertsResponse, err := ac.AlertService.GetAlerts(ctx.QueryParams())
	return ac.StandardResponse(ctx, alertsResponse, err)
}

func (ac *AlertsController) getAlert(ctx echo.Context) error {
	alertResponse, err := ac.AlertService.GetAlert(ctx.Param("alert"))
	return ac.StandardResponse(ctx, alertResponse, err)
}

func (ac *AlertsController) deleteAlert(ctx echo.Context) error {
	err := ac.AlertService.DeleteAlert(ctx.Param("alert"))
	return ac.StandardResponse(ctx, struct {
		Status string `json:"status"`
	}{Status: "ok"}, err)
}

func (ac *AlertsController) getAlertsCount(ctx echo.Context) error {
	alertsCount, err := ac.AlertService.GetAlertsCount(ctx.QueryParams())
	return ac.StandardResponse(ctx, alertsCount, err)
}

func (ac *AlertsController) getAlertsServices(ctx echo.Context) error {
	groupedServices, err := ac.AlertService.GetGroupedServices(ctx.QueryParams())
	return ac.StandardResponse(ctx, groupedServices, err)
}

func (ac *AlertsController) getAlertsEnvironments(ctx echo.Context) error {
	groupedEnvironments, err := ac.AlertService.GetGroupedEnvironments(ctx.QueryParams())
	return ac.StandardResponse(ctx, groupedEnvironments, err)
}

func (ac *AlertsController) updateAlertStatus(ctx echo.Context) error {
	var alertStatusUpdateRequest models.AlertStatusUpdateRequest
	err := ctx.Bind(&alertStatusUpdateRequest)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error()))
	}

	err = ac.AlertService.UpdateAlertStatus(ctx.Param("alert"), alertStatusUpdateRequest)
	return ac.StandardResponse(ctx, models.OK_RESPONSE, err)
}

func (ac *AlertsController) StandardResponse(ctx echo.Context, response interface{}, err error) error {
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse(err.Error()))
	}
	return ctx.JSON(http.StatusOK, response)
}
