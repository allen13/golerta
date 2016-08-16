package db

import (
	"github.com/allen13/golerta/app/models"
	"github.com/valyala/fasthttp"
)

type DB interface {
	Init() error
	CreateAlert(alert models.Alert) (string, error)
	CreateAlerts(alerts []models.Alert) (ids []string, err error)
	GetAlert(id string) (alert models.Alert, err error)
	DeleteAlert(id string) error
	UpdateAlert(id string, updates map[string]interface{}) error
	UpdateAlertStatus(id, status, text string)(err error)
	UpdateExistingAlertWithDuplicate(existingAlert models.Alert, duplicateAlert models.Alert) (err error)
	UpdateExistingAlertWithCorrelated(existingAlert models.Alert, correlatedAlert models.Alert) (err error)
	FindAlerts(queryArgs *fasthttp.Args) (alerts []models.Alert, err error)
	FindDuplicateAlert(alert models.Alert) (existingAlert models.Alert, alertIsDuplicate bool, err error)
	FindCorrelatedAlert(alert models.Alert) (existingAlert models.Alert, alertIsCorrelated bool, err error)
	CountAlertsGroup(group string, queryArgs *fasthttp.Args) (alertCountGroup map[string]int, err error)
	GetAlertEnvironmentsGroupedByEnvironment(queryArgs *fasthttp.Args) (groupedEnvironments []models.GroupedEnvironment, err error)
	GetAlertServicesGroupedByEnvironment(queryArgs *fasthttp.Args) (groupedServices []models.GroupedService, err error)
	EscalateTimedOutAlerts() error
}
