package models

import "time"

type AlertsResponse struct {
	Status         string    `json:"status"`
	Total          int       `json:"total"`
	Alerts         []Alert   `json:"alerts"`
	Page           int       `json:"page"`
	PageSize       int       `json:"pageSize"`
	Pages          int       `json:"pages"`
	More           bool      `json:"more"`
	SeverityCounts int       `json:"severityCounts"`
	StatusCounts   int       `json:"statusCounts"`
	LastTime       time.Time `json:"lastTime"`
	AutoRefresh    bool      `json:"autoRefresh"`
}

func NewAlertsResponse(alerts []Alert) (ar AlertsResponse) {
	ar = AlertsResponse{}
	ar.Alerts = alerts
	ar.Total = len(alerts)
	ar.Status = "ok"
	ar.AutoRefresh = true
	if ar.Total > 0{
		ar.LastTime = alerts[0].CreateTime
	}
	ar.More = false
	ar.Page = 1
	ar.Pages = 1
	ar.PageSize = 10000

	return
}
