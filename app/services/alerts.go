package services

import (
  "github.com/allen13/golerta/app/models"
  "time"
  "github.com/allen13/golerta/app/db"
)

type AlertService struct {
  db db.DB
}
func (as *AlertService) CreateAlert(alert models.Alert){
  alert.GenerateDefaults()
  as.db.CreateAlert(alert)
}

type HistoryEvent struct{
  Id string                     `gorethink:"id,omitempty" json:"id"`
  Event string                  `gorethink:"event" json:"event"`
  Severity string               `gorethink:"severtiy" json:"severtiy"`
  Value string                  `gorethink:"value" json:"value"`
  Type string                   `gorethink:"type" json:"type"`
  Text string                   `gorethink:"text" json:"text"`
  updateTime time.Time          `gorethink:"updateTime" json:"updateTime"`
}