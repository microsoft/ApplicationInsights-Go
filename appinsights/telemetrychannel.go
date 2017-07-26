package appinsights

import "time"

type TelemetryChannel interface {
	EndpointAddress() string
	Send(Telemetry)
	Flush()
	Stop()
	IsThrottled() bool
	Close(flush bool, retryTimeout time.Duration) chan bool
	CloseSync(flush bool, retryTimeout time.Duration)
}
