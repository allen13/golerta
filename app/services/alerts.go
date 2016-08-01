package services

import (
  "github.com/allen13/golerta/app/models"
  "github.com/allen13/golerta/app/db"
)

type AlertService struct {
  db db.DB
}

func (as *AlertService) ProcessAlert(alert models.Alert)(id string, err error){
  alert.GenerateDefaults()

  existingAlert, alertIsDuplicate, err := as.db.FindDuplicateAlert(alert)
  if err != nil {
    return
  }

  if alertIsDuplicate {
    err = as.db.UpdateExistingAlertWithDuplicate(existingAlert.Id, alert)
    if err != nil{
      return
    }

    id = existingAlert.Id
  }else {
    id, err = as.db.CreateAlert(alert)
  }

  return
}

func (as *AlertService) GetAlert(id string)(alert models.Alert, err error){
  alert, err = as.db.GetAlert(id)
  return
}

func (as *AlertService) DeleteAlert(id string)(err error){
  err = as.db.DeleteAlert(id)
  return
}

