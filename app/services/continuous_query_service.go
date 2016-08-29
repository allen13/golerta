package services

import (
	"github.com/allen13/golerta/app/algorithms"
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
	FlapDetection *algorithms.FlapDetection
}

func (cqs *ContinuousQueryService) Start() {
	go cqs.notifyPluginsOfOpenAlertsWithAlertableSeverity()

	queryTicker := time.NewTicker(cqs.QueryInterval)
	defer queryTicker.Stop()

	for {
		select {
		case <-queryTicker.C:
			go cqs.escalateTimedOutAlerts()
			go cqs.updateFlappingAlertScores()
		}
	}
}

func (cqs *ContinuousQueryService) updateFlappingAlertScores() {
	alerts, err := cqs.DB.FindFlappingAlerts()
	if err != nil {
		log.Println(err)
	}

	for _, alert := range alerts {
		isFlapping, currentFlapScore, remainingSeverityTimeChanges := cqs.FlapDetection.Detect(alert.SeverityChangeTimes)
		alert.FlapScore = currentFlapScore
		alert.SeverityChangeTimes = remainingSeverityTimeChanges
		err = cqs.DB.UpdateFlappingAlert(alert, isFlapping)
		if err != nil {
			log.Println(err)
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
