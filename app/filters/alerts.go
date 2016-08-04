package filters

import (
  "github.com/valyala/fasthttp"
  r "gopkg.in/dancannon/gorethink.v2"
)

func BuildAlertsFilter(queryArgs *fasthttp.Args)(rowFilter r.Term){
  rowFilter = r.Row
  for i,status := range getQueryValues("status", queryArgs){
    if i == 0{
      rowFilter = rowFilter.Field("status").Eq(status)
    } else {
      rowFilter = rowFilter.Or(r.Row.Field("status").Eq(status))
    }
  }
  return rowFilter
}

func getQueryValues(key string, queryArgs *fasthttp.Args)(values []string){
  for _,value := range queryArgs.PeekMulti(key){
    values = append(values, string(value))
  }
  return
}