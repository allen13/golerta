package notifiers

import (
	"github.com/allen13/golerta/app/models"
	"github.com/allen13/golerta/app/notifiers/file"
	"github.com/allen13/golerta/app/notifiers/pagerduty"
	"log"
)

type Notifiers struct {
	File              file.File           `toml:"file"`
	PagerDuty         pagerduty.PagerDuty `toml:"pagerduty"`
	TriggerSeverities []string            `toml:"trigger_severities"`
	notifiers         []Notifier
}

type Notifier interface {
	Trigger(alert models.Alert) error
	Acknowledge(alert models.Alert) error
	Resolve(alert models.Alert) error
	Init() error
	Enabled() bool
}

func (ns *Notifiers) Init() {
	if len(ns.TriggerSeverities) == 0 {
		ns.TriggerSeverities = []string{"critical"}
	}

	uninitializedNotifiers := []Notifier{&ns.File, &ns.PagerDuty}

	for _, notifier := range uninitializedNotifiers {
		if notifier.Enabled() {
			err := notifier.Init()
			if err != nil {
				log.Println(err)
			} else {
				ns.notifiers = append(ns.notifiers, notifier)
			}
		}
	}
}

func (ns *Notifiers) TriggerAlert(alert models.Alert) {
	for _, notifier := range ns.notifiers {
		err := notifier.Trigger(alert)
		if err != nil {
			log.Println(err)
		}
	}
}

func (ns *Notifiers) AcknowledgeAlert(alert models.Alert) {
	for _, notifier := range ns.notifiers {
		err := notifier.Acknowledge(alert)
		if err != nil {
			log.Println(err)
		}
	}
}

func (ns *Notifiers) ResolveAlert(alert models.Alert) {
	for _, notifier := range ns.notifiers {
		err := notifier.Resolve(alert)
		if err != nil {
			log.Println(err)
		}
	}
}

func (ns *Notifiers) ProcessAlertChangeFeed(alertChangeFeed models.AlertChangeFeed) {
	switch {
	case ns.IsAlertTriggered(alertChangeFeed):
		ns.TriggerAlert(alertChangeFeed.NewVal)

	case ns.IsAlertAcknowledged(alertChangeFeed):
		ns.AcknowledgeAlert(alertChangeFeed.NewVal)

	case ns.IsAlertResolved(alertChangeFeed):
		ns.ResolveAlert(alertChangeFeed.NewVal)
	}
}

func (ns *Notifiers) IsAlertTriggered(alertChangeFeed models.AlertChangeFeed) bool {
	return alertChangeFeed.NewVal.Status == "open" &&
		ns.alertHasTriggerSeverity(alertChangeFeed.NewVal) &&
		(severityChanged(alertChangeFeed) || statusChanged(alertChangeFeed))
}

func (ns *Notifiers) IsAlertAcknowledged(alertChangeFeed models.AlertChangeFeed) bool {
	return alertChangeFeed.NewVal.Status == "ack" &&
		alertChangeFeed.OldVal.Status == "open"
}

func (ns *Notifiers) IsAlertResolved(alertChangeFeed models.AlertChangeFeed) bool {
	return ((alertChangeFeed.NewVal.Status == "closed" ||
		alertChangeFeed.NewVal.Status == "silenced") &&
		statusChanged(alertChangeFeed)) ||
		ns.isAlertResolvedOnSeverityChange(alertChangeFeed)
}

//Resolve any alert that has a severity that is not a trigger but previously was
func (ns *Notifiers) isAlertResolvedOnSeverityChange(alertChangeFeed models.AlertChangeFeed) bool {
	return !ns.alertHasTriggerSeverity(alertChangeFeed.NewVal) && ns.alertHasTriggerSeverity(alertChangeFeed.OldVal)
}

func statusChanged(alertChangeFeed models.AlertChangeFeed) bool {
	return alertChangeFeed.NewVal.Status != alertChangeFeed.OldVal.Status
}

func severityChanged(alertChangeFeed models.AlertChangeFeed) bool {
	return alertChangeFeed.NewVal.Severity != alertChangeFeed.OldVal.Severity
}

func (ns *Notifiers) alertHasTriggerSeverity(alert models.Alert) bool {
	for _, severity := range ns.TriggerSeverities {
		if alert.Severity == severity {
			return true
		}
	}
	return false
}
