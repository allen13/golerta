package services

import (
  "github.com/allen13/golerta/app/models"
  "github.com/allen13/golerta/app/db"
)

type AlertService struct {
  db db.DB
}

func (as *AlertService) CreateAlert(alert models.Alert)(err error){
  alert.GenerateDefaults()
  _, err = as.db.CreateAlert(alert)
  return
}