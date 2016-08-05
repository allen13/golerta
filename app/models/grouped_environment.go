package models

type GroupedEnvironment struct {
	Environment string `gorethink:"environment" json:"environment"`
	Count       int    `gorethink:"count" json:"count"`
}

type GroupedEnvironmentResponse struct{
	Status         string    `json:"status"`
	Total          int       `json:"total"`
	Environments   []GroupedEnvironment `json:"environments"`
}

func NewGroupedEnvironmentResponse(groupedEnvironments []GroupedEnvironment) (g GroupedEnvironmentResponse) {
	g = GroupedEnvironmentResponse{}
	g.Status = "ok"
	g.Environments = groupedEnvironments

	for _, env := range groupedEnvironments {
		g.Total += env.Count
	}

	return
}