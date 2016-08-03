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
