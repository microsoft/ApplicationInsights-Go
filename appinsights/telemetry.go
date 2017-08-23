package appinsights

import (
	"fmt"
	"time"

	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
)

type TelemetryData interface {
	EnvelopeName() string
	BaseType() string
	//Properties() map[string]string
	//Measurements() map[string]float64
}

type Telemetry interface {
	Time() time.Time
	Context() TelemetryContext
	TelemetryData() TelemetryData
}

type BaseTelemetry struct {
	Timestamp time.Time
	context   TelemetryContext
}

func (item *BaseTelemetry) Time() time.Time {
	return item.Timestamp
}

func (item *BaseTelemetry) Context() TelemetryContext {
	return item.context
}

type TraceTelemetry struct {
	BaseTelemetry
	Data *contracts.MessageData
}

func NewTraceTelemetry(message string, severityLevel contracts.SeverityLevel) *TraceTelemetry {
	data := contracts.NewMessageData()
	data.Message = message

	return &TraceTelemetry{
		Data: data,
		BaseTelemetry: BaseTelemetry{
			Timestamp: time.Now(),
			context:   NewItemTelemetryContext(),
		},
	}
}

func (item *TraceTelemetry) TelemetryData() TelemetryData {
	return item.Data
}

type EventTelemetry struct {
	BaseTelemetry
	Data *contracts.EventData
}

func NewEventTelemetry(name string) *EventTelemetry {
	data := contracts.NewEventData()
	data.Name = name

	return &EventTelemetry{
		Data: data,
		BaseTelemetry: BaseTelemetry{
			Timestamp: time.Now(),
			context:   NewItemTelemetryContext(),
		},
	}
}

func (item *EventTelemetry) TelemetryData() TelemetryData {
	return item.Data
}

type MetricTelemetry struct {
	BaseTelemetry
	Data *contracts.MetricData
}

func NewMetricTelemetry(name string, value float64) *MetricTelemetry {
	dataPoint := contracts.NewDataPoint()
	dataPoint.Name = name
	dataPoint.Value = value
	dataPoint.Count = 1
	
	data := contracts.NewMetricData()
	data.Metrics = []*contracts.DataPoint{dataPoint}

	return &MetricTelemetry{
		Data: data,
		BaseTelemetry: BaseTelemetry{
			Timestamp: time.Now(),
			context:   NewItemTelemetryContext(),
		},
	}
}

func (item *MetricTelemetry) TelemetryData() TelemetryData {
	return item.Data
}

type RequestTelemetry struct {
	BaseTelemetry
	Data *contracts.RequestData
}

func NewRequestTelemetry(name, httpMethod, url string, timestamp time.Time, duration time.Duration, responseCode string, success bool) *RequestTelemetry {
	data := contracts.NewRequestData()
	data.Name = name
	data.StartTime = timestamp.Format(time.RFC3339Nano)
	data.Duration = formatDuration(duration)
	data.ResponseCode = responseCode
	data.Success = success
	data.HttpMethod = httpMethod
	data.Url = url
	data.Id = randomId()

	return &RequestTelemetry{
		Data: data,
		BaseTelemetry: BaseTelemetry{
			Timestamp: time.Now(),
			context:   NewItemTelemetryContext(),
		},
	}
}

func (item *RequestTelemetry) TelemetryData() TelemetryData {
	return item.Data
}

func formatDuration(d time.Duration) string {
	ticks := int64(d/(time.Nanosecond*100)) % 10000000
	seconds := int64(d/time.Second) % 60
	minutes := int64(d/time.Minute) % 60
	hours := int64(d/time.Hour) % 24
	days := int64(d / (time.Hour * 24))

	return fmt.Sprintf("%d.%02d:%02d:%02d.%07d", days, hours, minutes, seconds, ticks)
}
