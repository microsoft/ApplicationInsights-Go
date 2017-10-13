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

func (trace *TraceTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewMessageData()
	data.Message = trace.Message
	data.Properties = trace.Properties
	data.SeverityLevel = trace.SeverityLevel

	return data
}

type EventTelemetry struct {
	BaseTelemetry
	Name string
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

func (event *EventTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewEventData()
	data.Name = event.Name
	data.Properties = event.Properties
	data.Measurements = event.Measurements

	return data
}

type MetricTelemetry struct {
	BaseTelemetry
	Name  string
	Value float64
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

type MetricsTelemetry struct {
	BaseTelemetry
	Metrics []*contracts.DataPoint
}

func NewMetricsTelemetry() *MetricsTelemetry {
	return &MetricsTelemetry{
		BaseTelemetry: BaseTelemetry{
			Timestamp:  currentClock.Now(),
			Context:    NewTelemetryContext(),
			Properties: make(map[string]string),
		},
	}
}

func (metrics *MetricsTelemetry) AddMeasurement(name string, value float64) {
	dataPoint := contracts.NewDataPoint()
	dataPoint.Name = name
	dataPoint.Kind = Measurement
	dataPoint.Value = value
	dataPoint.Count = 1

	metrics.Metrics = append(metrics.Metrics, dataPoint)
}

func (metrics *MetricsTelemetry) AddAggregation(name string, value, min, max, stddev float64, count int) {
	dataPoint := contracts.NewDataPoint()
	dataPoint.Name = name
	dataPoint.Kind = Aggregation
	dataPoint.Value = value
	dataPoint.Min = min
	dataPoint.Max = max
	dataPoint.StdDev = stddev
	dataPoint.Count = count

	metrics.Metrics = append(metrics.Metrics, dataPoint)
}

func (metrics *MetricsTelemetry) AddMeasurements(measurements map[string]float64) {
	for k, v := range measurements {
		dataPoint := contracts.NewDataPoint()
		dataPoint.Name = k
		dataPoint.Kind = Measurement
		dataPoint.Value = v
		dataPoint.Count = 1

		metrics.Metrics = append(metrics.Metrics, dataPoint)
	}
}

func (metrics *MetricsTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewMetricData()
	data.Metrics = metrics.Metrics
	data.Properties = metrics.Properties
	return data
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

func (request *RequestTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewRequestData()
	data.Name = request.Name
	data.Duration = formatDuration(request.Duration)
	data.ResponseCode = request.ResponseCode
	data.Success = request.Success
	data.Url = request.Url

	if request.Id == "" {
		data.Id = RandomId()
	} else {
		data.Id = request.Id
	}

	data.Properties = request.Properties
	data.Measurements = request.Measurements
	return data
}

type RemoteDependencyTelemetry struct {
	BaseTelemetry
	Name       string
	Id         string
	ResultCode string
	Duration   time.Duration
	Success    bool
	Data       string
	Type       string
	Target     string
}

func NewRemoteDependencyTelemetry(dependencyType, target string, success bool) *RemoteDependencyTelemetry {
	return &RemoteDependencyTelemetry{
		Type:    dependencyType,
		Target:  target,
		Success: success,
		BaseTelemetry: BaseTelemetry{
			Timestamp:    currentClock.Now(),
			Context:      NewTelemetryContext(),
			Properties:   make(map[string]string),
			Measurements: make(map[string]float64),
		},
	}
}

func (telem *RemoteDependencyTelemetry) MarkTime(startTime, endTime time.Time) {
	telem.Timestamp = startTime
	telem.Duration = endTime.Sub(startTime)
}

func (telem *RemoteDependencyTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewRemoteDependencyData()
	data.Name = telem.Name
	data.Id = telem.Id
	data.ResultCode = telem.ResultCode
	data.Duration = formatDuration(telem.Duration)
	data.Success = telem.Success
	data.Data = telem.Data
	data.Target = telem.Target
	data.Properties = telem.Properties
	data.Measurements = telem.Measurements

	return data
}

type AvailabilityTelemetry struct {
	BaseTelemetry
	Id          string
	Name        string
	Duration    time.Duration
	Success     bool
	RunLocation string
	Message     string
}

func NewAvailabilityTelemetry(name string, duration time.Duration, success bool) *AvailabilityTelemetry {
	return &AvailabilityTelemetry{
		Name:     name,
		Duration: duration,
		Success:  success,
		Id:       RandomId(),
		BaseTelemetry: BaseTelemetry{
			Timestamp:    currentClock.Now(),
			Context:      NewTelemetryContext(),
			Properties:   make(map[string]string),
			Measurements: make(map[string]float64),
		},
	}
}

func (telem *AvailabilityTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewAvailabilityData()
	data.Id = telem.Id
	data.Name = telem.Name
	data.Duration = formatDuration(telem.Duration)
	data.Success = telem.Success
	data.RunLocation = telem.RunLocation
	data.Message = telem.Message
	data.Properties = telem.Properties
	data.Measurements = telem.Measurements
	return data
}

type PageViewTelemetry struct {
	BaseTelemetry
	Url      string
	Duration time.Duration
	Name     string
}

func NewPageViewTelemetry(url string) *PageViewTelemetry {
	return &PageViewTelemetry{
		Url: url,
		BaseTelemetry: BaseTelemetry{
			Timestamp:    currentClock.Now(),
			Context:      NewTelemetryContext(),
			Properties:   make(map[string]string),
			Measurements: make(map[string]float64),
		},
	}
}

func (telem *PageViewTelemetry) TelemetryData() TelemetryData {
	data := contracts.NewPageViewData()
	data.Url = telem.Url
	data.Duration = formatDuration(telem.Duration)
	data.Name = telem.Name
	data.Properties = telem.Properties
	data.Measurements = telem.Measurements
	return data
}

func formatDuration(d time.Duration) string {
	ticks := int64(d/(time.Nanosecond*100)) % 10000000
	seconds := int64(d/time.Second) % 60
	minutes := int64(d/time.Minute) % 60
	hours := int64(d/time.Hour) % 24
	days := int64(d / (time.Hour * 24))

	return fmt.Sprintf("%d.%02d:%02d:%02d.%07d", days, hours, minutes, seconds, ticks)
}
