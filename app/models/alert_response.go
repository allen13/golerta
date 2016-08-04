package models

type AlertResponse struct {
	Status string `json:"status"`
	Total  int    `json:"total"`
	Alert  Alert  `json:"alert"`
}

func NewAlertResponse(alert Alert) (a AlertResponse) {
	a = AlertResponse{}
	a.Alert = alert
	a.Total = 1
	a.Status = "ok"

	return
}
