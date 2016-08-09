package services

import (
	"github.com/allen13/golerta/app/db"
	"github.com/allen13/golerta/app/models"
	"github.com/valyala/fasthttp"
)

type AlertService struct {
	DB db.DB
}

func (as *AlertService) ProcessAlert(alert models.Alert) (id string, err error) {
	alert.GenerateDefaults()

	//Check for duplicate alerts
	existingAlert, alertIsDuplicate, err := as.DB.FindDuplicateAlert(alert)
	if err != nil {
		return
	}

	if alertIsDuplicate {
		err = as.DB.UpdateExistingAlertWithDuplicate(existingAlert.Id, alert)
		if err != nil {
			return
		}

		id = existingAlert.Id
		return
	}

	//Check for correlated alerts
	existingCorrelatedAlert, alertIsCorrelated, err := as.DB.FindCorrelatedAlert(alert)
	if err != nil {
		return
	}

	if alertIsCorrelated {
		err = as.DB.UpdateExistingAlertWithCorrelated(existingCorrelatedAlert, alert)
		if err != nil {
			return
		}

		id = existingCorrelatedAlert.Id
		return
	}

	//Alert is neither duplicate or correlated, create a new one
	id, err = as.DB.CreateAlert(alert)

	return
}

func (as *AlertService) GetAlert(id string) (alertResponse models.AlertResponse, err error) {
	alert, err := as.DB.GetAlert(id)
	alertResponse = models.NewAlertResponse(alert)
	return
}

func (as *AlertService) GetAlerts(queryArgs *fasthttp.Args) (alertsResponse models.AlertsResponse, err error) {
	alerts, err := as.DB.FindAlerts(queryArgs)
	if err != nil {
		return
	}
	alertsResponse = models.NewAlertsResponse(alerts)

	return
}

func (as *AlertService) DeleteAlert(id string) (err error) {
	err = as.DB.DeleteAlert(id)
	return
}

func (as *AlertService) GetAlertsCount(queryArgs *fasthttp.Args)(models.AlertsCountResponse,error) {
	severityCounts, err := as.DB.CountAlertsGroup("severity", queryArgs)
	if err != nil{
		return models.AlertsCountResponse{}, err
	}
	statusCounts, err := as.DB.CountAlertsGroup("status", queryArgs)
	if err != nil{
		return models.AlertsCountResponse{}, err
	}

	return models.NewAlertsCountResponse(statusCounts, severityCounts), nil
}

func (as *AlertService) GetGroupedServices(queryArgs *fasthttp.Args)(models.GroupedServiceResponse,error) {
	groupedServices, err := as.DB.GetAlertServicesGroupedByEnvironment(queryArgs)
	if err != nil{
		return models.GroupedServiceResponse{}, err
	}

	return models.NewGroupedServiceResponse(groupedServices), nil
}

func (as *AlertService) GetGroupedEnvironments(queryArgs *fasthttp.Args)(models.GroupedEnvironmentResponse,error) {
	groupedEnvironments, err := as.DB.GetAlertEnvironmentsGroupedByEnvironment(queryArgs)
	if err != nil{
		return models.GroupedEnvironmentResponse{}, err
	}

	return models.NewGroupedEnvironmentResponse(groupedEnvironments), nil
}

func (as *AlertService) UpdateAlertStatus(id, status, text string) (err error){
	err = as.DB.UpdateAlertStatus(id, status, text)
	return
}
