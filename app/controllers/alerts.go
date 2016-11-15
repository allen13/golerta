package controllers

import (
	"github.com/allen13/golerta/app/models"
	"github.com/allen13/golerta/app/services"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
)

type AlertsController struct {
	Echo             *echo.Echo
	AlertService     services.AlertService
	AuthMiddleware   echo.MiddlewareFunc
	LogAlertRequests bool
}

func (ac *AlertsController) Init() {
	ac.Echo.POST("/alert", ac.createAlert, ac.AuthMiddleware)
	ac.Echo.GET("/alert/:alert", ac.getAlert, ac.AuthMiddleware)
	ac.Echo.GET("/alerts", ac.getAlerts, ac.AuthMiddleware)
	ac.Echo.POST("/alert/:alert/status", ac.updateAlertStatus, ac.AuthMiddleware)
	ac.Echo.DELETE("/alert/:alert", ac.deleteAlert, ac.AuthMiddleware)
	ac.Echo.GET("/alerts/count", ac.getAlertsCount, ac.AuthMiddleware)
	ac.Echo.GET("/alerts/services", ac.getAlertsServices, ac.AuthMiddleware)
	ac.Echo.GET("/alerts/environments", ac.getAlertsEnvironments, ac.AuthMiddleware)

}

func (ac *AlertsController) createAlert(ctx echo.Context) error {
	if ac.LogAlertRequests {
		request, _ := ioutil.ReadAll(ctx.Request().Body)
		log.Println(string(request))
	}

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
