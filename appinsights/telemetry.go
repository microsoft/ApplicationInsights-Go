package appinsights

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
)

type TelemetryData interface {
	EnvelopeName() string
	BaseType() string
}

type Telemetry interface {
	Time() time.Time
	TelemetryContext() *TelemetryContext
	TelemetryData() TelemetryData
	GetProperties() map[string]string
	GetMeasurements() map[string]float64
}

type BaseTelemetry struct {
	Timestamp    time.Time
	Properties   map[string]string
	Measurements map[string]float64
	Context      *TelemetryContext
}

func (item *BaseTelemetry) Time() time.Time {
	return item.Timestamp
}

func (item *BaseTelemetry) TelemetryContext() *TelemetryContext {
	return item.Context
}

func (item *BaseTelemetry) GetProperties() map[string]string {
	return item.Properties
}

func (item *BaseTelemetry) GetMeasurements() map[string]float64 {
	return item.Measurements
}

type TraceTelemetry struct {
	BaseTelemetry
	Message       string
	SeverityLevel contracts.SeverityLevel
}

func (trace *TraceTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewMessageData()
	data.Message = trace.Message
	data.Properties = trace.Properties
	data.SeverityLevel = trace.SeverityLevel

	return data
}

func NewTraceTelemetry(message string, severityLevel contracts.SeverityLevel) *TraceTelemetry {
	return &TraceTelemetry{
		Message:       message,
		SeverityLevel: severityLevel,
		BaseTelemetry: BaseTelemetry{
			Timestamp:  currentClock.Now(),
			Context:    NewTelemetryContext(),
			Properties: make(map[string]string),
		},
	}
}

type EventTelemetry struct {
	BaseTelemetry
	Name string
}

func (event *EventTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewEventData()
	data.Name = event.Name
	data.Properties = event.Properties
	data.Measurements = event.Measurements

	return data
}

func NewEventTelemetry(name string) *EventTelemetry {
	return &EventTelemetry{
		Name: name,
		BaseTelemetry: BaseTelemetry{
			Timestamp:    currentClock.Now(),
			Context:      NewTelemetryContext(),
			Properties:   make(map[string]string),
			Measurements: make(map[string]float64),
		},
	}
}

type MetricTelemetry struct {
	BaseTelemetry
	Name  string
	Value float64
}

func (metric *MetricTelemetry) TelemetryData() TelemetryData {
	dataPoint := contracts.NewDataPoint()
	dataPoint.Name = metric.Name
	dataPoint.Value = metric.Value
	dataPoint.Count = 1

	data := contracts.NewMetricData()
	data.Metrics = []*contracts.DataPoint{dataPoint}
	data.Properties = metric.Properties

	return data
}

func NewMetricTelemetry(name string, value float64) *MetricTelemetry {
	return &MetricTelemetry{
		Name:  name,
		Value: value,
		BaseTelemetry: BaseTelemetry{
			Timestamp:  currentClock.Now(),
			Context:    NewTelemetryContext(),
			Properties: make(map[string]string),
		},
	}
}

type MetricsTelemetry struct {
	BaseTelemetry
	Metrics []*contracts.DataPoint
}

func (metrics *MetricsTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewMetricData()
	data.Metrics = metrics.Metrics
	data.Properties = metrics.Properties
	return data
}

func NewMetricsTelemetry(values map[string]float64) *MetricsTelemetry {
	var dataPoints []*contracts.DataPoint
	for k, v := range values {
		dataPoint := contracts.NewDataPoint()
		dataPoint.Name = k
		dataPoint.Value = v
		dataPoint.Count = 1

		dataPoints = append(dataPoints, dataPoint)
	}

	return &MetricsTelemetry{
		Metrics: dataPoints,
		BaseTelemetry: BaseTelemetry{
			Timestamp:  currentClock.Now(),
			Context:    NewTelemetryContext(),
			Properties: make(map[string]string),
		},
	}
}

type RequestTelemetry struct {
	BaseTelemetry
	Id           string
	Name         string
	Url          string
	Duration     time.Duration
	ResponseCode string
	Success      bool
}

func (request *RequestTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewRequestData()
	data.Name = request.Name
	data.Duration = formatDuration(request.Duration)
	data.ResponseCode = request.ResponseCode
	data.Success = request.Success
	data.Url = request.Url

	if request.Id == "" {
		data.Id = RandomId()
	}

	data.Properties = request.Properties
	data.Measurements = request.Measurements
	return data
}

func NewRequestTelemetry(method, url string, duration time.Duration, responseCode string) *RequestTelemetry {
	success := true
	code, err := strconv.Atoi(responseCode)
	if err != nil {
		success = code < 400 || code == 401
	}

	return &RequestTelemetry{
		Name:         fmt.Sprintf("%s %s", method, url),
		Url:          url,
		Id:           RandomId(),
		Duration:     duration,
		ResponseCode: responseCode,
		Success:      success,
		BaseTelemetry: BaseTelemetry{
			Timestamp:    currentClock.Now().Add(-duration),
			Context:      NewTelemetryContext(),
			Properties:   make(map[string]string),
			Measurements: make(map[string]float64),
		},
	}
}

func formatDuration(d time.Duration) string {
	ticks := int64(d/(time.Nanosecond*100)) % 10000000
	seconds := int64(d/time.Second) % 60
	minutes := int64(d/time.Minute) % 60
	hours := int64(d/time.Hour) % 24
	days := int64(d / (time.Hour * 24))

	return fmt.Sprintf("%d.%02d:%02d:%02d.%07d", days, hours, minutes, seconds, ticks)
}
