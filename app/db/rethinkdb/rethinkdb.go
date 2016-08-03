package rethinkdb

import (
  r "gopkg.in/dancannon/gorethink.v2"
  "github.com/allen13/golerta/app/models"
)

type RethinkDB struct {
  Address string
  Database string
  session* r.Session
}

func (re* RethinkDB) Init()(error){
  if re.Address == "" {
    re.Address = "localhost:28015"
  }
  if re.Database == "" {
    re.Database = "alerta"
  }

  return re.Connect()
}
func (re* RethinkDB) Connect()(error){
  session, err := r.Connect(r.ConnectOpts{
    Address: re.Address,
  })
  if err != nil {
    return err
  }
  re.session = session

  err = re.CreateDBIfNotExist()
  if err != nil {
    return err
  }

  err = re.CreateTableIfNotExist("alerts")
  if err != nil {
    return err
  }

  return nil
}

func (re* RethinkDB) CreateDBIfNotExist()(error){
  exists, err := re.DBExists()
  if err != nil {
    return err
  }

  if !exists{
    _, err := r.DBCreate(re.Database).RunWrite(re.session)
    if err != nil {
      return err
    }
  }

  return nil
}

func (re* RethinkDB) CreateTableIfNotExist(table string)(error){
  exists, err := re.TableExists(table)
  if err != nil {
    return err
  }

  if !exists{
    _, err := r.DB(re.Database).TableCreate(table).RunWrite(re.session)
    if err != nil {
      return err
    }
  }

  return nil
}

func (re* RethinkDB) DBExists()(bool, error){
  var response []interface{}
  res, err := r.DBList().Run(re.session)

  if err != nil{
    return false, err
  }

  err = res.All(&response)
  if err != nil{
    return false, err
  }

  for _, db := range response{
    if db == re.Database{
      return true, nil
    }
  }

  return false, nil
}

func (re* RethinkDB) TableExists(table string)(bool, error){
  var response []interface{}
  res, err := r.DB(re.Database).TableList().Run(re.session)

  if err != nil{
    return false, err
  }

  err = res.All(&response)
  if err != nil{
    return false, err
  }

  for _, responseTable := range response{
    if responseTable == table{
      return true, nil
    }
  }

  return false, nil
}

//Create alert and return generated id
func (re* RethinkDB) CreateAlert(alert models.Alert)(string, error){
  ids, err := re.CreateAlerts([]models.Alert{alert})
  if err != nil{
    return "", err
  }
  if len(ids) < 1{
    return alert.Id, nil
  }
  return ids[0], nil
}

//Create alerts and return generated ids
func (re* RethinkDB) CreateAlerts(alerts []models.Alert)(ids []string, err error){
  writeResponse, err := r.DB(re.Database).Table("alerts").Insert(alerts).RunWrite(re.session)
  if err != nil {
    return ids, err
  }
  return writeResponse.GeneratedKeys, nil
}

func (re* RethinkDB) FindOneAlert(filter interface{})(alert models.Alert, foundOne bool, err error) {
  alerts, err := re.FindAlerts(filter)
  if err != nil {
    return
  }
  if len(alerts) < 1{
    return
  }
  if len(alerts) >= 1{
    foundOne = true
    alert = alerts[0]
    return
  }
  return
}

func (re* RethinkDB) FindAlerts(filter interface{})(alerts []models.Alert, err error) {
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

func (re* RethinkDB) FindDuplicateAlert(alert models.Alert)(existingAlert models.Alert, alertIsDuplicate bool, err error) {
  findDuplicateAlert := map[string]interface{}{
    "event": alert.Event,
    "environment": alert.Environment,
    "resource": alert.Resource,
    "severity": alert.Severity,
    "customer": alert.Customer,
  }

  existingAlert, alertIsDuplicate, err = re.FindOneAlert(findDuplicateAlert)

  return
}

func (re* RethinkDB) FindCorrelatedAlert(alert models.Alert)(existingAlert models.Alert, alertIsCorrelated bool, err error) {
  var correlatedFilter = func(user r.Term) r.Term {
    return user.Field("event").Eq(alert.Event).And(
           user.Field("environment").Eq(alert.Environment)).And(
           user.Field("resource").Eq(alert.Resource)).And(
           user.Field("customer").Eq(alert.Customer)).And(
           user.Field("severity").Ne(alert.Severity))
  }

  existingAlert, alertIsCorrelated, err = re.FindOneAlert(correlatedFilter)

  return
}



func (re* RethinkDB) GetAlert(id string)(alert models.Alert, err error) {
  res, err := r.DB(re.Database).Table("alerts").Get(id).Run(re.session)
  if err != nil {
    return
  }
  defer res.Close()
  err = res.One(&alert)
  return
}

func (re* RethinkDB) DeleteAlert(id string)(error) {
  _, err := r.DB(re.Database).Table("alerts").Get(id).Delete().RunWrite(re.session)
  if err != nil {
    return err
  }
  return nil
}

func (re* RethinkDB) UpdateAlert(id string, updates map[string]interface{})(error) {
  _, err := r.DB(re.Database).Table("alerts").Get(id).Update(updates).RunWrite(re.session)
  if err != nil {
    return err
  }
  return nil
}

func (re* RethinkDB) UpdateExistingAlertWithDuplicate(existingId string, duplicateAlert models.Alert)(err error){
  alertUpdate := map[string]interface{}{
    "value": duplicateAlert.Value,
    "text": duplicateAlert.Text,
    "tags": duplicateAlert.Tags,
    "rawData": duplicateAlert.RawData,
    "repeat": true,
    "lastReceiveId": duplicateAlert.Id,
    "lastReceiveTime": duplicateAlert.ReceiveTime,
    "duplicateCount": r.Row.Field("duplicateCount").Add(1),
    "history": r.Row.Field("history").Prepend(duplicateAlert.History[0]),
  }
  err = re.UpdateAlert(existingId, alertUpdate)
  return
}

func (re* RethinkDB) UpdateExistingAlertWithCorrelated(existingAlert models.Alert, correlatedAlert models.Alert)(err error){
  alertUpdate := map[string]interface{}{
    "severity": correlatedAlert.Severity,
    "previousSeverity": existingAlert.Severity,
    "status": correlatedAlert.Status,
    "value": correlatedAlert.Value,
    "text": correlatedAlert.Text,
    "tags": correlatedAlert.Tags,
    "createTime": correlatedAlert.CreateTime,
    "rawData": correlatedAlert.RawData,
    "duplicateCount": 0,
    "repeat": false,
    "lastReceiveId": correlatedAlert.Id,
    "lastReceiveTime": correlatedAlert.ReceiveTime,
    "history": r.Row.Field("history").Prepend(correlatedAlert.History[0]),
  }
  err = re.UpdateAlert(existingAlert.Id, alertUpdate)
  return
}
