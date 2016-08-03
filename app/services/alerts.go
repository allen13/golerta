package services

import (
  "github.com/allen13/golerta/app/models"
  "github.com/allen13/golerta/app/db"
)

type AlertService struct {
  DB db.DB
}

func (as *AlertService) ProcessAlert(alert models.Alert)(id string, err error){
  alert.GenerateDefaults()

  //Check for duplicate alerts
  existingAlert, alertIsDuplicate, err := as.DB.FindDuplicateAlert(alert)
  if err != nil {
    return
  }

  if alertIsDuplicate {
    err = as.DB.UpdateExistingAlertWithDuplicate(existingAlert.Id, alert)
    if err != nil{
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
    if err != nil{
      return
    }

    id = existingCorrelatedAlert.Id
    return
  }

  //Alert is neither duplicate or correlated, create a new one
  id, err = as.DB.CreateAlert(alert)

  return
}

func (as *AlertService) GetAlert(id string)(alert models.Alert, err error){
  alert, err = as.DB.GetAlert(id)
  return
}

func (as *AlertService) GetAlerts()(alertResponse models.AlertsResponse, err error){
  alertResponse = models.AlertsResponse{}
  emptyFilter := map[string]string{}
  alertResponse.Alerts, err = as.DB.FindAlerts(emptyFilter)
  if err != nil {
    return
  }
  alertResponse.Status = "ok"
  alertResponse.Total = len(alertResponse.Alerts)
  return
}

func (as *AlertService) DeleteAlert(id string)(err error){
  err = as.DB.DeleteAlert(id)
  return
}

