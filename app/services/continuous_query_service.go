package services

import (
	"time"
	"github.com/allen13/golerta/app/db/rethinkdb"
	"log"
)

type ContinuousQueryService struct {
	DB rethinkdb.RethinkDB
	QueryInterval time.Duration
}

func (cqs *ContinuousQueryService) Start(){
	queryTicker := time.NewTicker(cqs.QueryInterval)
	defer queryTicker.Stop()

	for {
		select {
		case <-queryTicker.C:
			go cqs.EscalateTimedOutAlerts()
		}
	}
}

func (cqs *ContinuousQueryService) EscalateTimedOutAlerts(){
	err := cqs.DB.EscalateTimedOutAlerts()
	if err != nil{
		log.Println(err)
	}
}