package services

import (
	"github.com/allen13/golerta/app/algorithms"
	"github.com/allen13/golerta/app/db/rethinkdb"
	"github.com/allen13/golerta/app/models"
	"github.com/valyala/fasthttp"
	"log"
)

type AlertService struct {
	DB            *rethinkdb.RethinkDB
	FlapDetection *algorithms.FlapDetection
}

func (as *AlertService) ProcessAlert(currentAlert models.Alert) (id string, err error) {
	currentAlert.GenerateDefaults()
	existingRelatedAlert, foundExistingRelatedAlert, err := as.DB.FindRelatedAlert(currentAlert)

	if !foundExistingRelatedAlert {
		//New Alert
		id, err = as.DB.CreateAlert(currentAlert)
		if err != nil {
			log.Println(err)
		}
		return
	}

	alertSeverityStateChanged := as.detectSeverityChange(&currentAlert, existingRelatedAlert)

	if !alertSeverityStateChanged {
		//Duplicate Alert

		err = as.DB.UpdateExistingAlertWithDuplicate(existingRelatedAlert, currentAlert)
		if err != nil {
			log.Println(err)
		}

		id = existingRelatedAlert.Id

		return

	} else {
		//Correlated Alert

		err = as.DB.UpdateExistingAlertWithCorrelated(existingRelatedAlert, currentAlert)
		if err != nil {
			log.Println(err)
		}

		id = existingRelatedAlert.Id
		return
	}

}

func (as *AlertService) detectSeverityChange(currentAlert *models.Alert, existingRelatedAlert models.Alert) (alertSeverityStateChanged bool) {
	alertSeverityStateChanged = existingRelatedAlert.Severity != currentAlert.Severity

	if !as.FlapDetection.Enabled {
		return
	}

	previouslyFlapping := existingRelatedAlert.Severity == "flapping"
	flapSeverityStateChanged := existingRelatedAlert.FlapSeverityState != currentAlert.Severity

	alertSeverityStateChanged = (previouslyFlapping && flapSeverityStateChanged) ||
		(!previouslyFlapping && alertSeverityStateChanged)

	if alertSeverityStateChanged {
		//Add new severity change time if the state changed
		existingRelatedAlert.SeverityChangeTimes = append(existingRelatedAlert.SeverityChangeTimes, currentAlert.CreateTime)
	}

	isFlapping, flapScore, remainingSeverityTimeChanges := as.FlapDetection.Detect(existingRelatedAlert.SeverityChangeTimes)
	currentAlert.SeverityChangeTimes = remainingSeverityTimeChanges
	currentAlert.FlapScore = flapScore

	switch {
	case isFlapping:
		currentAlert.FlapSeverityState = currentAlert.Severity
		currentAlert.Severity = "flapping"
		alertSeverityStateChanged = !previouslyFlapping
	case !isFlapping && previouslyFlapping:
		alertSeverityStateChanged = true
	}

	return
}

func (as *AlertService) GetAlert(id string) (alertResponse models.AlertResponse, err error) {
	alert, err := as.DB.GetAlert(id)
	alertResponse = models.NewAlertResponse(alert)
	return
}

func (as *AlertService) GetAlerts(queryArgs *fasthttp.Args) (alertsResponse models.AlertsResponse, err error) {
	alerts, err := as.DB.GetAlertsSummary(queryArgs)
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

func (as *AlertService) GetAlertsCount(queryArgs *fasthttp.Args) (models.AlertsCountResponse, error) {
	severityCounts, err := as.DB.CountAlertsGroup("severity", queryArgs)
	if err != nil {
		return models.AlertsCountResponse{}, err
	}
	statusCounts, err := as.DB.CountAlertsGroup("status", queryArgs)
	if err != nil {
		return models.AlertsCountResponse{}, err
	}

	return models.NewAlertsCountResponse(statusCounts, severityCounts), nil
}

func (as *AlertService) GetGroupedServices(queryArgs *fasthttp.Args) (models.GroupedServiceResponse, error) {
	groupedServices, err := as.DB.GetAlertServicesGroupedByEnvironment(queryArgs)
	if err != nil {
		return models.GroupedServiceResponse{}, err
	}

	return models.NewGroupedServiceResponse(groupedServices), nil
}

func (as *AlertService) GetGroupedEnvironments(queryArgs *fasthttp.Args) (models.GroupedEnvironmentResponse, error) {
	groupedEnvironments, err := as.DB.GetAlertEnvironmentsGroupedByEnvironment(queryArgs)
	if err != nil {
		return models.GroupedEnvironmentResponse{}, err
	}

	return models.NewGroupedEnvironmentResponse(groupedEnvironments), nil
}

func (as *AlertService) UpdateAlertStatus(id string, alertStatusUpdateRequest models.AlertStatusUpdateRequest) (err error) {
	err = as.DB.UpdateAlertStatus(id, alertStatusUpdateRequest.Status, alertStatusUpdateRequest.Text, alertStatusUpdateRequest.AcknowledgementDuration)
	return
}
