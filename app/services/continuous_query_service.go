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
	go cqs.processStreamingAlertChanges()

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

func (cqs *ContinuousQueryService) processStreamingAlertChanges() {
	alertsChannel := make(chan models.AlertChangeFeed)

	err := cqs.DB.StreamAlertChanges(alertsChannel)
	if err != nil {
		close(alertsChannel)
	}

CHANGE_FEED_LOOP:
	for {
		select {
		case alertChangeFeed, ok := <-alertsChannel:
			if !ok {
				log.Println("Alerts change feed closed.")
				break CHANGE_FEED_LOOP
			}
			//Send alerts to notifier plugins
			go cqs.Notifiers.ProcessAlertChangeFeed(alertChangeFeed)
		}
	}

	log.Println("Processing alert change feed failed. Trying again in 10 seconds...")
	time.Sleep(time.Second * 10)
	go cqs.processStreamingAlertChanges()
}
