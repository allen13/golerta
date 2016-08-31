package models

type GroupedService struct {
	Service     string `gorethink:"service" json:"service"`
	Environment string `gorethink:"environment" json:"environment"`
	Count       int    `gorethink:"count" json:"count"`
}

type GroupedServiceResponse struct {
	Status   string           `json:"status"`
	Total    int              `json:"total"`
	Services []GroupedService `json:"services"`
}

func NewGroupedServiceResponse(groupedServices []GroupedService) (g GroupedServiceResponse) {
	g = GroupedServiceResponse{}
	g.Status = "ok"
	g.Services = groupedServices

	for _, env := range groupedServices {
		g.Total += env.Count
	}

	return
}
