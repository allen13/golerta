package db

import "github.com/allen13/golerta/app/models"

type DB interface {
	Connect() error
	CreateDBIfNotExist() error
	CreateTableIfNotExist(table string) error
	DBExists() (bool, error)
	TableExists(table string) (bool, error)
	CreateAlert(alert models.Alert) (string, error)
	CreateAlerts(alerts []models.Alert) (ids []string, err error)
	GetAlert(id string) (alert models.Alert, err error)
	DeleteAlert(id string) error
	UpdateAlert(id string, updates map[string]interface{}) error
	UpdateExistingAlertWithDuplicate(existingId string, duplicateAlert models.Alert) (err error)
	UpdateExistingAlertWithCorrelated(existingAlert models.Alert, correlatedAlert models.Alert) (err error)
	FindAlerts(filter interface{}) (alerts []models.Alert, err error)
	FindOneAlert(filter interface{}) (alert models.Alert, foundOne bool, err error)
	FindDuplicateAlert(alert models.Alert) (existingAlert models.Alert, alertIsDuplicate bool, err error)
	FindCorrelatedAlert(alert models.Alert) (existingAlert models.Alert, alertIsCorrelated bool, err error)
	CountAlertsGroup(group string, filter interface{}) (alertCountGroup map[string]int, err error)
	GetAlertEnvironmentsGroupedByEnvironment(filter interface{}) (groupedEnvironments []models.GroupedEnvironment, err error)
	GetAlertServicesGroupedByEnvironment(filter interface{}) (groupedServices []models.GroupedService, err error)
}
