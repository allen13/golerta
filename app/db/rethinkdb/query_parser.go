package rethinkdb

import (
	r "gopkg.in/dancannon/gorethink.v2"
)

var ALLOWED_QUERY_PARAMS []string = []string{
	"id",
	"status",
	"environment",
	"service",
}

func BuildAlertsFilter(queryParams map[string][]string) (rowFilter r.Term) {
	if len(queryParams) < 1 {
		return r.Row
	}

	var firstParam = true
	for _, allowedQueryParam := range ALLOWED_QUERY_PARAMS {
		_, hasParam := queryParams[allowedQueryParam]
		if !hasParam {
			continue
		}

		paramFilter := buildQueryForParam(allowedQueryParam, queryParams)

		if firstParam {
			rowFilter = paramFilter
			firstParam = false
		} else {
			rowFilter = rowFilter.And(paramFilter)
		}
	}

	return rowFilter
}

func buildQueryForParam(queryParam string, queryParamValues map[string][]string) r.Term {
	paramFilter := r.Row

	for i, queryValue := range queryParamValues[queryParam] {
		if queryParam == "service" {
			if i == 0 {
				paramFilter = paramFilter.Field(queryParam).Contains(queryValue)
			} else {
				paramFilter = paramFilter.Or(r.Row.Field(queryParam).Contains(queryValue))
			}
		} else {
			if i == 0 {
				paramFilter = paramFilter.Field(queryParam).Eq(queryValue)
			} else {
				paramFilter = paramFilter.Or(r.Row.Field(queryParam).Eq(queryValue))
			}
		}

	}

	return paramFilter
}
