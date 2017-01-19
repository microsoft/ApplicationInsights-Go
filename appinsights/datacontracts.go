package appinsights

import (
	"fmt"
	"time"
)

type Telemetry struct {
	Timestamp time.Time
	Context   TelemetryContext
	TypeName  string
	Data      interface{}
}

func NewTraceTelemetry(message string, severityLevel SeverityLevel) Telemetry {
	return Telemetry{
		Timestamp: time.Now(),
		Context:   NewItemTelemetryContext(),
		TypeName:  "Message",
		Data: &MessageData{
			Message: message,
			Ver:     2}}
}

func NewEventTelemetry(name string) Telemetry {
	return Telemetry{
		Timestamp: time.Now(),
		Context:   NewItemTelemetryContext(),
		TypeName:  "Event",
		Data: &EventData{
			Name: name,
			Ver:  2}}
}

func NewMetricTelemetry(name string, value float32) Telemetry {
	return Telemetry{
		Timestamp: time.Now(),
		Context:   NewItemTelemetryContext(),
		TypeName:  "Metric",
		Data: &MetricData{
			Metrics: []*DataPoint{&DataPoint{
				Name:  name,
				Value: value,
				Count: 1}},
			Ver: 2}}
}

func NewRequestTelemetry(
	id string,
	name string,
	timestamp time.Time,
	duration time.Duration,
	httpMethod string,
	url string,
	responseCode string,
	success bool) Telemetry {

	return Telemetry{
		Timestamp: time.Now(),
		Context:   NewItemTelemetryContext(),
		TypeName:  "Request",
		Data: &RequestData{
			Id:           id,
			Name:         name,
			StartTime:    timestamp.Format(time.RFC3339Nano),
			Duration:     formatDuration(duration),
			ResponseCode: responseCode,
			Success:      success,
			HttpMethod:   httpMethod,
			Url:          url,
			Ver:          2}}
}

func NewRemoteDependencyData(
	id string,
	name string,
	resultCode int,
	commandName string,
	kind DataPointType,
	duration time.Duration,
	count int,
	min float32,
	max float32,
	stdDev float32,
	theType string,
	dependencyKind DependencyKind,
	success bool,
	async bool,
	dependencySource DependencySourceType,
	properties map[string]string,
	alterTelementryContext func(*TelemetryContext)) Telemetry {

	context := NewItemTelemetryContext()
	alterTelementryContext(&context)

	return Telemetry{
		Timestamp: time.Now(),
		Context:   context,
		TypeName:  "RemoteDependency",
		Data: &RemoteDependencyData{
			Ver:              2,
			Id:               id,
			Name:             name,
			ResultCode:       resultCode,
			CommandName:      commandName,
			Kind:             kind,
			Duration:         formatDuration(duration),
			Count:            count,
			Min:              min,
			Max:              max,
			StdDev:           stdDev,
			Type:             theType,
			DependencyKind:   dependencyKind,
			Success:          success,
			Async:            async,
			DependencySource: dependencySource,
			Properties:       properties}}
}

func formatDuration(duration time.Duration) string {
	var (
		refHours   = int(duration.Hours())
		refMinutes = int(duration.Minutes())
		refSeconds = int(duration.Seconds())
		refMilli   = int(duration.Nanoseconds() / 1e3)
		days       = refHours / 24
		hours      = refHours - days*24
		minutes    = refMinutes - refHours*60
		seconds    = refSeconds - refMinutes*60
		milli      = refMilli - refSeconds*1e6
	)
	return fmt.Sprintf("%02d.%02d:%02d:%02d.%04d",
		days,
		hours,
		minutes,
		seconds,
		milli)
}
