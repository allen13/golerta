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

func (re* RethinkDB) GetAllAlerts(filter map[string]interface{})(alerts []models.Alert, err error) {
  res, err := r.DB(re.Database).Table("alerts").Filter(filter).Run(re.session)
  if err != nil {
    return
  }
  defer res.Close()
  err = res.All(&alerts)

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
