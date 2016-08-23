package services

import (
	"github.com/allen13/golerta/app/db/rethinkdb"
	"github.com/allen13/golerta/app/models"
	"github.com/allen13/golerta/app/notifiers"
	"log"
	"time"
)

type ContinuousQueryService struct {
	DB            rethinkdb.RethinkDB
	QueryInterval time.Duration
	Notifiers     notifiers.Notifiers
}

func (cqs *ContinuousQueryService) Start() {
	queryTicker := time.NewTicker(cqs.QueryInterval)
	defer queryTicker.Stop()

	go cqs.notifyPluginsOfOpenAlertsWithAlertableSeverity()

	for {
		select {
		case <-queryTicker.C:
			go cqs.escalateTimedOutAlerts()
		}
	}
}

func (cqs *ContinuousQueryService) escalateTimedOutAlerts() {
	err := cqs.DB.EscalateTimedOutAlerts()
	if err != nil {
		log.Println(err)
	}
}

func (cqs *ContinuousQueryService) notifyPluginsOfOpenAlertsWithAlertableSeverity() {
	alertsChannel := make(chan models.AlertChangeFeed)

	cqs.DB.StreamAlertChanges(alertsChannel)
	for {
		select {
		case alertChangeFeed := <-alertsChannel:
			cqs.Notifiers.ProcessAlertChangeFeed(alertChangeFeed)
		}
	}
}
