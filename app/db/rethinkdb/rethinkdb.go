package rethinkdb

import (
	"github.com/allen13/golerta/app/models"
	r "gopkg.in/dancannon/gorethink.v2"
	"github.com/valyala/fasthttp"
)

type RethinkDB struct {
	Address  string
	Database string
	session  *r.Session
}

func (re *RethinkDB) Init() error {
	if re.Address == "" {
		re.Address = "localhost:28015"
	}
	if re.Database == "" {
		re.Database = "alerta"
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

func (re *RethinkDB) FindAlerts(queryArgs *fasthttp.Args) (alerts []models.Alert, err error) {
	filter := BuildAlertsFilter(queryArgs)
	return re.findAlerts(filter)
}

func (re *RethinkDB) FindDuplicateAlert(alert models.Alert) (existingAlert models.Alert, alertIsDuplicate bool, err error) {
	findDuplicateAlert := map[string]interface{}{
		"event":       alert.Event,
		"environment": alert.Environment,
		"resource":    alert.Resource,
		"severity":    alert.Severity,
		"customer":    alert.Customer,
	}

	existingAlert, alertIsDuplicate, err = re.findOneAlert(findDuplicateAlert)

	return
}

func (re *RethinkDB) FindCorrelatedAlert(alert models.Alert) (existingAlert models.Alert, alertIsCorrelated bool, err error) {
	var correlatedFilter = func(user r.Term) r.Term {
		return user.Field("event").Eq(alert.Event).And(
			user.Field("environment").Eq(alert.Environment)).And(
			user.Field("resource").Eq(alert.Resource)).And(
			user.Field("customer").Eq(alert.Customer)).And(
			user.Field("severity").Ne(alert.Severity))
	}

	existingAlert, alertIsCorrelated, err = re.findOneAlert(correlatedFilter)

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

func (re *RethinkDB) UpdateExistingAlertWithDuplicate(existingId string, duplicateAlert models.Alert) (err error) {
	alertUpdate := map[string]interface{}{
		"value":           duplicateAlert.Value,
		"text":            duplicateAlert.Text,
		"tags":            duplicateAlert.Tags,
		"rawData":         duplicateAlert.RawData,
		"repeat":          true,
		"lastReceiveId":   duplicateAlert.Id,
		"lastReceiveTime": duplicateAlert.ReceiveTime,
		"duplicateCount":  r.Row.Field("duplicateCount").Add(1),
		"history":         r.Row.Field("history").Prepend(duplicateAlert.History[0]),
	}
	err = re.UpdateAlert(existingId, alertUpdate)
	return
}

func (re *RethinkDB) UpdateExistingAlertWithCorrelated(existingAlert models.Alert, correlatedAlert models.Alert) (err error) {
	alertUpdate := map[string]interface{}{
		"severity":         correlatedAlert.Severity,
		"previousSeverity": existingAlert.Severity,
		"status":           correlatedAlert.Status,
		"value":            correlatedAlert.Value,
		"text":             correlatedAlert.Text,
		"tags":             correlatedAlert.Tags,
		"createTime":       correlatedAlert.CreateTime,
		"rawData":          correlatedAlert.RawData,
		"duplicateCount":   0,
		"repeat":           false,
		"lastReceiveId":    correlatedAlert.Id,
		"lastReceiveTime":  correlatedAlert.ReceiveTime,
		"history":          r.Row.Field("history").Prepend(correlatedAlert.History[0]),
	}
	err = re.UpdateAlert(existingAlert.Id, alertUpdate)
	return
}

func (re *RethinkDB) GetAlertServicesGroupedByEnvironment(queryArgs *fasthttp.Args) (groupedServices []models.GroupedService, err error) {
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

func (re *RethinkDB) GetAlertEnvironmentsGroupedByEnvironment(queryArgs *fasthttp.Args) (groupedEnvironments []models.GroupedEnvironment, err error) {
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

func (re *RethinkDB) CountAlertsGroup(group string, queryArgs *fasthttp.Args) (alertCountGroup map[string]int, err error) {
	filter := BuildAlertsFilter(queryArgs)
	t := r.DB(re.Database).Table("alerts").Filter(filter).Group(group).Count().Ungroup().Map(
		r.Object(r.Row.Field("group"), r.Row.Field("reduction"))).Reduce(func(left r.Term, right r.Term)(r.Term){
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