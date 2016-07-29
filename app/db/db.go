package db

import "github.com/allen13/golerta/app/models"

type DB interface{
  Connect()(error)
  CreateDBIfNotExist()(error)
  CreateTableIfNotExist(table string)(error)
  DBExists()(bool, error)
  TableExists(table string)(bool, error)
  CreateAlert(alert models.Alert)(string, error)
  CreateAlerts(alerts []models.Alert)(ids []string, err error)
  GetAllAlerts(filter map[string]interface{})(alerts []models.Alert, err error)
  GetAlert(id string)(alert models.Alert, err error)
  DeleteAlert(id string)(error)
  UpdateAlert(id string, updates map[string]interface{})(error)
}

