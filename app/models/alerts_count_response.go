package models

type AlertsCountResponse struct {
	Status         string         `json:"status"`
	Total          int            `json:"total"`
	StatusCounts   map[string]int `json:"statusCounts"`
	SeverityCounts map[string]int `json:"severityCounts"`
}

func NewAlertsCountResponse(statusCounts, severityCounts map[string]int) (a AlertsCountResponse) {
	a = AlertsCountResponse{}
	a.Status = "ok"
	a.SeverityCounts = severityCounts
	a.StatusCounts = statusCounts

	for _, count := range severityCounts {
		a.Total += count
	}

	return
}
