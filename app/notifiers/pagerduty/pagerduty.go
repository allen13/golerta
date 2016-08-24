package pagerduty

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/allen13/golerta/app/models"
)

type PagerDuty struct {
	EnabledField bool   `toml:"enabled"`
	ServiceKey   string `toml:"service_key"`
}

func (pd *PagerDuty) Init() error {
	return nil
}

func (pd *PagerDuty) Enabled() bool {
	return pd.EnabledField
}

func (pd *PagerDuty) CreatePagerDutyEvent(eventType string, alert models.Alert) error {
	event := pagerduty.Event{
		ServiceKey:  pd.ServiceKey,
		Type:        eventType,
		Description: alert.String(),
		IncidentKey: alert.Id,
	}

	_, err := pagerduty.CreateEvent(event)
	return err
}

func (pd *PagerDuty) Trigger(alert models.Alert) error {
	return pd.CreatePagerDutyEvent("trigger", alert)
}

func (pd *PagerDuty) Acknowledge(alert models.Alert) error {
	return pd.CreatePagerDutyEvent("acknowledge", alert)
}

func (pd *PagerDuty) Resolve(alert models.Alert) error {
	return pd.CreatePagerDutyEvent("resolve", alert)
}
