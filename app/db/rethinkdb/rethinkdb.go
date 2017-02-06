package rethinkdb

import (
	"time"

	"github.com/allen13/golerta/app/models"
	r "gopkg.in/dancannon/gorethink.v2"
)

type RethinkDB struct {
	Address           string `toml:"address"`
	Database          string `toml:"database"`
	AlertHistoryLimit int    `toml:"alert_history_limit"`
	session           *r.Session
}

func (re *RethinkDB) Init() error {
	if re.Address == "" {
		re.Address = "localhost:28015"
	}
	if re.Database == "" {
		re.Database = "alerta"
	}
	if re.AlertHistoryLimit == 0 {
		re.AlertHistoryLimit = 100
	}

	return re.connect()
}
func (re *RethinkDB) connect() error {
	session, err := r.Connect(r.ConnectOpts{
		Address: re.Address,
	})
	if err != nil {
		return err
	}
	re.session = session

	err = re.createDBIfNotExist()
	if err != nil {
		return err
	}

	err = re.createTableIfNotExist("alerts")
	if err != nil {
		return err
	}

	return nil
}

func (re *RethinkDB) createDBIfNotExist() error {
	exists, err := re.dbExists()
	if err != nil {
		return err
	}

	if !exists {
		_, err := r.DBCreate(re.Database).RunWrite(re.session)
		if err != nil {
			return err
		}
	}

	return nil
}

func (re *RethinkDB) createTableIfNotExist(table string) error {
	exists, err := re.tableExists(table)
	if err != nil {
		return err
	}

	if !exists {
		_, err := r.DB(re.Database).TableCreate(table).RunWrite(re.session)
		if err != nil {
			return err
		}
	}

	return nil
}

func (re *RethinkDB) dbExists() (bool, error) {
	var response []interface{}
	res, err := r.DBList().Run(re.session)

	if err != nil {
		return false, err
	}

	err = res.All(&response)
	if err != nil {
		return false, err
	}

	for _, db := range response {
		if db == re.Database {
			return true, nil
		}
	}

	return false, nil
}

func (re *RethinkDB) tableExists(table string) (bool, error) {
	var response []interface{}
	res, err := r.DB(re.Database).TableList().Run(re.session)

	if err != nil {
		return false, err
	}

	err = res.All(&response)
	if err != nil {
		return false, err
	}

	for _, responseTable := range response {
		if responseTable == table {
			return true, nil
		}
	}

	return false, nil
}

//Create alert and return generated id
func (re *RethinkDB) CreateAlert(alert models.Alert) (string, error) {
	ids, err := re.CreateAlerts([]models.Alert{alert})
	if err != nil {
		return "", err
	}
	if len(ids) < 1 {
		return alert.Id, nil
	}
	return ids[0], nil
}

//Create alerts and return generated ids
func (re *RethinkDB) CreateAlerts(alerts []models.Alert) (ids []string, err error) {
	writeResponse, err := r.DB(re.Database).Table("alerts").Insert(alerts).RunWrite(re.session)
	if err != nil {
		return ids, err
	}
	return writeResponse.GeneratedKeys, nil
}

func (re *RethinkDB) findOneAlert(filter interface{}) (alert models.Alert, foundOne bool, err error) {
	alerts, err := re.findAlerts(filter)
	if err != nil {
		return
	}
	if len(alerts) < 1 {
		return
	}
	if len(alerts) >= 1 {
		foundOne = true
		alert = alerts[0]
		return
	}
	return
}

func (re *RethinkDB) findAlerts(filter interface{}) (alerts []models.Alert, err error) {
	res, err := r.DB(re.Database).Table("alerts").Filter(filter).Run(re.session)
	if err != nil {
		return
	}
	defer res.Close()
	err = res.All(&alerts)
	if alerts == nil {
		alerts = []models.Alert{}
	}
	return
}

func (re *RethinkDB) FindFlappingAlerts() ([]models.Alert, error) {
	flappingAlerts := map[string]interface{}{
		"severity": "flapping",
	}
	return re.findAlerts(flappingAlerts)
}

func (re *RethinkDB) FindAlerts(queryArgs map[string][]string) (alerts []models.Alert, err error) {
	filter := BuildAlertsFilter(queryArgs)
	return re.findAlerts(filter)
}

//A filtered list of alerts without unecessary details
func (re *RethinkDB) GetAlertsSummary(queryArgs map[string][]string) (alertsSummary []map[string]interface{}, err error) {
	filter := BuildAlertsFilter(queryArgs)
	res, err := r.DB(re.Database).Table("alerts").Filter(filter).WithFields(
		"id",
		"severity",
		"status",
		"lastReceiveTime",
		"environment",
		"service",
		"resource",
		"event",
		"value").Run(re.session)
	if err != nil {
		return
	}
	defer res.Close()
	err = res.All(&alertsSummary)
	if alertsSummary == nil {
		alertsSummary = make([]map[string]interface{}, 0)
	}
	return
}

func (re *RethinkDB) FindRelatedAlert(alert models.Alert) (relatedAlert models.Alert, foundAlert bool, err error) {
	findRelatedAlert := map[string]interface{}{
		"event":       alert.Event,
		"environment": alert.Environment,
		"resource":    alert.Resource,
		"customer":    alert.Customer,
	}

	relatedAlert, foundAlert, err = re.findOneAlert(findRelatedAlert)

	return
}

func (re *RethinkDB) GetAlert(id string) (alert models.Alert, err error) {
	res, err := r.DB(re.Database).Table("alerts").Get(id).Run(re.session)
	if err != nil {
		return
	}
	defer res.Close()
	err = res.One(&alert)
	return
}

func (re *RethinkDB) DeleteAlert(id string) error {
	_, err := r.DB(re.Database).Table("alerts").Get(id).Delete().RunWrite(re.session)
	if err != nil {
		return err
	}
	return nil
}

func (re *RethinkDB) UpdateAlert(id string, updates map[string]interface{}) error {
	_, err := r.DB(re.Database).Table("alerts").Get(id).Update(updates).RunWrite(re.session)
	if err != nil {
		return err
	}
	return nil
}

func (re *RethinkDB) UpdateExistingAlertWithDuplicate(existingAlert models.Alert, duplicateAlert models.Alert) (err error) {
	alertUpdate := map[string]interface{}{
		"value":               duplicateAlert.Value,
		"text":                duplicateAlert.Text,
		"tags":                duplicateAlert.Tags,
		"rawData":             duplicateAlert.RawData,
		"repeat":              true,
		"lastReceiveId":       duplicateAlert.Id,
		"lastReceiveTime":     duplicateAlert.ReceiveTime,
		"duplicateCount":      r.Row.Field("duplicateCount").Add(1),
		"timeout":             duplicateAlert.Timeout,
		"flapScore":           duplicateAlert.FlapScore,
		"severityChangeTimes": duplicateAlert.SeverityChangeTimes,
		"flapSeverityState":   duplicateAlert.FlapSeverityState,
	}

	if existingAlert.Status == "resolved" {
		alertUpdate["status"] = "open"
	}

	if existingAlert.Status != duplicateAlert.Status {
		alertUpdate["history"] = r.Row.Field("history").Limit(re.AlertHistoryLimit).Prepend(r.Object(
			"id", duplicateAlert.Id,
			"status", duplicateAlert.Status,
			"event", duplicateAlert.Event,
			"value", duplicateAlert.Value,
			"type", "duplicate alert update",
			"updateTime", duplicateAlert.CreateTime,
		))
	}

	err = re.UpdateAlert(existingAlert.Id, alertUpdate)
	return
}

func (re *RethinkDB) UpdateExistingAlertWithCorrelated(existingAlert models.Alert, correlatedAlert models.Alert) (err error) {
	alertUpdate := map[string]interface{}{
		"severity":            correlatedAlert.Severity,
		"previousSeverity":    existingAlert.Severity,
		"value":               correlatedAlert.Value,
		"text":                correlatedAlert.Text,
		"tags":                correlatedAlert.Tags,
		"createTime":          correlatedAlert.CreateTime,
		"rawData":             correlatedAlert.RawData,
		"duplicateCount":      0,
		"repeat":              false,
		"lastReceiveId":       correlatedAlert.Id,
		"lastReceiveTime":     correlatedAlert.ReceiveTime,
		"timeout":             correlatedAlert.Timeout,
		"flapScore":           correlatedAlert.FlapScore,
		"flapSeverityState":   correlatedAlert.FlapSeverityState,
		"severityChangeTimes": correlatedAlert.SeverityChangeTimes,
		"history": r.Row.Field("history").Limit(re.AlertHistoryLimit).Prepend(r.Object(
			"id", correlatedAlert.Id,
			"severity", correlatedAlert.Severity,
			"event", correlatedAlert.Event,
			"value", correlatedAlert.Value,
			"type", "correlated alert update",
			"updateTime", correlatedAlert.CreateTime,
		)),
	}

	if existingAlert.Status == "resolved" {
		alertUpdate["status"] = "open"
	}

	err = re.UpdateAlert(existingAlert.Id, alertUpdate)
	return
}

func (re *RethinkDB) UpdateFlappingAlert(alert models.Alert, isFlapping bool) (err error) {
	alertUpdate := map[string]interface{}{
		"flapScore":           alert.FlapScore,
		"severityChangeTimes": alert.SeverityChangeTimes,
	}

	if !isFlapping {
		alertUpdate["severity"] = alert.FlapSeverityState
		alertUpdate["flapSeverityState"] = ""
		alertUpdate["history"] = r.Row.Field("history").Limit(re.AlertHistoryLimit).Prepend(r.Object(
			"id", alert.Id,
			"severity", alert.FlapSeverityState,
			"event", alert.Event,
			"value", alert.Value,
			"type", "continuous flap check update",
			"updateTime", time.Now(),
		))
	}

	err = re.UpdateAlert(alert.Id, alertUpdate)
	return
}

func (re *RethinkDB) GetAlertServicesGroupedByEnvironment(queryArgs map[string][]string) (groupedServices []models.GroupedService, err error) {
	filter := BuildAlertsFilter(queryArgs)

	//This query unwinds the service field which duplicates the alert per service. It then groups by both environment and service.
	//Most of the time the service field only has one service but the alerta api made this field into a list for some reason.
	//If it is decided that the service field can only represent a single service then this query can be greatly simplified.
	t := r.DB(re.Database).Table("alerts").Filter(filter).ConcatMap(func(alert r.Term) r.Term {
		return alert.Field("service").Map(func(service r.Term) r.Term {
			return r.Object(
				"service", service,
				"environment", alert.Field("environment"))
		})
	}).CoerceTo("array").Group("environment", "service").Count().Ungroup().Map(func(result r.Term) r.Term {
		return r.Object(
			"environment", result.Field("group").AtIndex(0),
			"service", result.Field("group").AtIndex(1),
			"count", result.Field("reduction"))
	})

	res, err := t.Run(re.session)
	if err != nil {
		return
	}
	defer res.Close()
	err = res.All(&groupedServices)
	if groupedServices == nil {
		groupedServices = []models.GroupedService{}
	}
	return
}

func (re *RethinkDB) GetAlertEnvironmentsGroupedByEnvironment(queryArgs map[string][]string) (groupedEnvironments []models.GroupedEnvironment, err error) {
	filter := BuildAlertsFilter(queryArgs)
	t := r.DB(re.Database).Table("alerts").Filter(filter).Group("environment").Count().Ungroup().Map(func(result r.Term) r.Term {
		return r.Object(
			"environment", result.Field("group"),
			"count", result.Field("reduction"))
	})

	res, err := t.Run(re.session)
	if err != nil {
		return
	}
	defer res.Close()
	err = res.All(&groupedEnvironments)
	if groupedEnvironments == nil {
		groupedEnvironments = []models.GroupedEnvironment{}
	}
	return
}

func (re *RethinkDB) CountAlertsGroup(group string, queryArgs map[string][]string) (alertCountGroup map[string]int, err error) {
	filter := BuildAlertsFilter(queryArgs)
	t := r.DB(re.Database).Table("alerts").Filter(filter).Group(group).Count().Ungroup().Map(
		r.Object(r.Row.Field("group"), r.Row.Field("reduction"))).Reduce(func(left r.Term, right r.Term) r.Term {
		return left.Merge(right)
	})

	res, err := t.Run(re.session)
	if err != nil {
		return
	}
	defer res.Close()
	err = res.One(&alertCountGroup)
	if alertCountGroup == nil {
		alertCountGroup = map[string]int{}
	}
	return
}

func (re *RethinkDB) UpdateAlertStatus(id, status, text string, acknowledgementDuration int) (err error) {
	alert, err := re.GetAlert(id)
	if err != nil {
		return err
	}

	historyEvent := models.HistoryEvent{
		Id:         alert.Id,
		Event:      alert.Event,
		Status:     status,
		Value:      alert.Value,
		Text:       text,
		Type:       "update by user",
		UpdateTime: time.Now(),
	}

	updates := map[string]interface{}{
		"status":  status,
		"history": r.Row.Field("history").Limit(re.AlertHistoryLimit).Prepend(historyEvent),
	}

	if status == "ack" {
		updates["acknowledgement_duration"] = acknowledgementDuration
		updates["acknowledgementTime"] = time.Now()
	}

	err = re.UpdateAlert(id, updates)
	return
}

func (re *RethinkDB) EscalateTimedOutAlerts() error {
	timedOutAlerts := r.Row.Field("severity").Ne("critical").And(
		r.Row.Field("timeout").Ne(0)).And(
		r.Row.Field("lastReceiveTime").Add(r.Row.Field("timeout")).Lt(r.Now()))

	criticalUpdate := map[string]interface{}{
		"severity": "critical",
		"value":    "ALERT TIMED OUT",
		"history": r.Row.Field("history").Limit(re.AlertHistoryLimit).Prepend(r.Object(
			"id", r.Row.Field("id"),
			"severity", "critical",
			"event", r.Row.Field("event"),
			"value", "ALERT TIMED OUT",
			"type", "continuous query - time out",
			"updateTime", time.Now(),
		)),
	}

	updateEachTimedOutAlert := r.DB(re.Database).Table("alerts").Filter(timedOutAlerts).Update(criticalUpdate)

	_, err := updateEachTimedOutAlert.RunWrite(re.session)
	if err != nil {
		return err
	}
	return nil
}

func (re *RethinkDB) ReopenAwknowledgedAlers() error {
	expiredAcknowledgedAlerts := r.Row.Field("status").Eq("ack").And(
		r.Row.Field("acknowledgement_duration").Ne(0)).And(
		r.Row.Field("acknowledgementTime").Add(r.Row.Field("acknowledgement_duration")).Lt(r.Now()))

	statusUpdate := map[string]interface{}{
		"status":                   "open",
		"acknowledgement_duration": 0,
		"history": r.Row.Field("history").Limit(re.AlertHistoryLimit).Prepend(r.Object(
			"id", r.Row.Field("id"),
			"status", "open",
			"event", r.Row.Field("event"),
			"value", r.Row.Field("value"),
			"type", "continuous query - reopen acknowledged alert",
			"updateTime", time.Now(),
		)),
	}

	updateExpiredAcknowledgedAlerts := r.DB(re.Database).Table("alerts").Filter(expiredAcknowledgedAlerts).Update(statusUpdate)

	_, err := updateExpiredAcknowledgedAlerts.RunWrite(re.session)
	if err != nil {
		return err
	}
	return nil
}

func (re *RethinkDB) StreamAlertChanges(alertsChannel chan models.AlertChangeFeed) (err error) {
	changesOpts := r.ChangesOpts{
		IncludeTypes:   true,
		IncludeInitial: false,
	}

	cursor, err := r.DB(re.Database).Table("alerts").Changes(changesOpts).Run(re.session)
	if err != nil {
		return
	}

	cursor.Listen(alertsChannel)
	return
}
