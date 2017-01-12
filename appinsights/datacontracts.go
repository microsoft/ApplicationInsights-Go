package appinsights

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
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
	name string,
	timestamp time.Time,
	duration time.Duration,
	responseCode string,
	success bool) Telemetry {

	return Telemetry{
		Timestamp: time.Now(),
		Context:   NewItemTelemetryContext(),
		TypeName:  "Request",
		Data: &RequestData{
			Id:           uuid.NewV4().String(),
			Name:         name,
			StartTime:    timestamp.Format(time.RFC3339Nano),
			Duration:     formatDuration(duration),
			ResponseCode: responseCode,
			Success:      success,
			Ver:          2}}
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
