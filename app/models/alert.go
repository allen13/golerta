package models

import (
  "time"
  "github.com/twinj/uuid"
)

type Alert struct {
  Id string                     `gorethink:"id,omitempty" json:"id"`
  Event string                  `gorethink:"event" json:"event"`
  Severity string               `gorethink:"severity" json:"severity"`
  DuplicateCount int            `gorethink:"duplicateCount" json:"duplicateCount"`
  Resource string               `gorethink:"resource" json:"resource"`
  Environment string            `gorethink:"environment" json:"environment"`
  Correlate []string            `gorethink:"correlate" json:"correlate"`
  Status string                 `gorethink:"status" json:"status"`
  Service []string              `gorethink:"service" json:"service"`
  Group string                  `gorethink:"group" json:"group"`
  Value string                  `gorethink:"value" json:"value"`
  Text string                   `gorethink:"text" json:"text"`
  Tags []string                 `gorethink:"tags" json:"tags"`
  Attributes map[string]string  `gorethink:"attributes" json:"attributes"`
  Origin string                 `gorethink:"origin" json:"origin"`
  EventType string              `gorethink:"eventType" json:"eventType"`
  CreateTime time.Time          `gorethink:"createTime" json:"createTime"`
  Timeout int                   `gorethink:"timeout" json:"timeout"`
  RawData string                `gorethink:"rawData" json:"rawData"`
  Customer string               `gorethink:"customer" json:"customer"`
  Repeat bool                   `gorethink:"repeat" json:"repeat"`
  PreviousSeverity string       `gorethink:"previousSeverity" json:"previousSeverity"`
  TrendIndication string        `gorethink:"trendIndication" json:"trendIndication"`
  ReceiveTime time.Time         `gorethink:"receiveTime" json:"receiveTime"`
  LastReceiveId string          `gorethink:"lastReceiveId" json:"lastReceiveId"`
  LastReceiveTime time.Time     `gorethink:"lastReceiveTime" json:"lastReceiveTime"`
  History []HistoryEvent        `gorethink:"history" json:"history"`
}

type HistoryEvent struct{
  Id string                     `gorethink:"id,omitempty" json:"id"`
  Event string                  `gorethink:"event" json:"event"`
  Severity string               `gorethink:"severity" json:"severity"`
  Value string                  `gorethink:"value" json:"value"`
  Type string                   `gorethink:"type" json:"type"`
  Text string                   `gorethink:"text" json:"text"`
  UpdateTime time.Time          `gorethink:"updateTime" json:"updateTime"`
}

func (alert* Alert) GenerateDefaults(){
  if alert.Id == "" {
    id := uuid.NewV4()
    alert.Id = id.String()
  }

  if alert.Attributes == nil{
    alert.Attributes = map[string]string{}
  }

  if alert.CreateTime.IsZero() {
    alert.CreateTime = time.Now()
  }

  if alert.ReceiveTime.IsZero() {
    alert.ReceiveTime = time.Now()
  }

  if alert.LastReceiveTime.IsZero() {
    alert.LastReceiveTime = time.Now()
  }

  alert.History = []HistoryEvent{HistoryEvent{
    Id: alert.Id,
    Event: alert.Event,
    Severity: alert.Severity,
    Value: alert.Value,
    Type: alert.EventType,
    Text: alert.Text,
    UpdateTime: alert.CreateTime,
  }}
}