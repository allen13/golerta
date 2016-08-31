package services

import (
	"github.com/allen13/golerta/app/algorithms"
	"github.com/allen13/golerta/app/db/rethinkdb"
	"github.com/allen13/golerta/app/models"
	"testing"
)

func TestAlertService_ProcessAlert(t *testing.T) {
	db := getTestDB(t)
	as := &AlertService{
		db,
		&algorithms.FlapDetection{},
	}

	alert := models.Alert{
		Event:       "duplicate alert",
		Resource:    "testServer01",
		Environment: "syd01",
		Severity:    "critical",
		Origin:      "consul-syd01",
	}

	alertId1, err := as.ProcessAlert(alert)
	if err != nil {
		t.Fatal(err)
	}

	alertId2, err := as.ProcessAlert(alert)
	if err != nil {
		t.Fatal(err)
	}

	if alertId1 != alertId2 {
		err = as.DeleteAlert(alertId1)
		if err != nil {
			t.Fatal(err)
		}
		err = as.DeleteAlert(alertId2)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatal("Created two separate alerts instead of updating the first duplicate alert")
	}

	retrievedAlertResponse, err := as.GetAlert(alertId1)
	if err != nil {
		t.Fatal(err)
	}

	if retrievedAlertResponse.Alert.DuplicateCount < 1 {
		t.Fatal("Failed to update duplicate field")
	}

	if retrievedAlertResponse.Alert.DuplicateCount > 1 {
		t.Fatal("Duplicate field was updated too many times")
	}

	err = as.DeleteAlert(alertId1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAlertService_ProcessCorrelatedAlerts(t *testing.T) {
	db := getTestDB(t)
	as := &AlertService{
		db,
		&algorithms.FlapDetection{},
	}

	alert := models.Alert{
		Event:       "correlated alert",
		Resource:    "testServer01",
		Environment: "syd01",
		Severity:    "critical",
		Origin:      "consul-syd01",
	}

	correlatedAlert := models.Alert{
		Event:       "correlated alert",
		Resource:    "testServer01",
		Environment: "syd01",
		Severity:    "normal",
		Origin:      "consul-syd01",
	}

	alertId1, err := as.ProcessAlert(alert)
	if err != nil {
		t.Fatal(err)
	}

	alertId2, err := as.ProcessAlert(correlatedAlert)
	if err != nil {
		t.Fatal(err)
	}

	if alertId1 != alertId2 {
		err = as.DeleteAlert(alertId1)
		if err != nil {
			t.Fatal(err)
		}
		err = as.DeleteAlert(alertId2)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatal("Created two separate alerts instead of updating the first correlated alert")
	}

	retrievedAlertResponse, err := as.GetAlert(alertId1)
	if err != nil {
		t.Fatal(err)
	}

	if retrievedAlertResponse.Alert.DuplicateCount != 0 {
		err = as.DeleteAlert(alertId1)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatal("Failed to clear duplicate count")
	}

	err = as.DeleteAlert(alertId1)
	if err != nil {
		t.Fatal(err)
	}
}

//docker run -d --name rethinkdb -p 8080:8080 -p 28015:28015 rethinkdb
func getTestDB(t *testing.T) (db *rethinkdb.RethinkDB) {
	db = &rethinkdb.RethinkDB{}
	err := db.Init()

	if err != nil {
		t.Fatal(err)
	}

	return db
}
