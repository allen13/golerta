package services

import (
  "testing"
  "github.com/allen13/golerta/app/db/rethinkdb"
  "github.com/allen13/golerta/app/db"
  "github.com/allen13/golerta/app/models"
)

func TestAlertService_CreateAlert(t *testing.T) {
  db := getTestDB(t)
  as := &AlertService{db}

  alert := models.Alert{
    Event: "cpu usage idle",
    Resource: "testServer01",
    Environment: "syd01",
    Severity: "CRITICAL",
    Origin: "consul-syd01",
  }

  err := as.CreateAlert(alert)
  if err != nil {
    t.Fatal(err)
  }
}

//docker run -d --name rethinkdb -p 8080:8080 -p 28015:28015 rethinkdb
func getTestDB(t *testing.T)(db.DB){
  db := &rethinkdb.RethinkDB{
    Address: "localhost:28015",
    Database: "alerta",
  }
  err := db.Connect()

  if err != nil{
    t.Fatal(err)
  }

  return db
}
