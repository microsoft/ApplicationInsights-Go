package appinsights

import "time"

type TelemetryChannel interface {
	EndpointAddress() string
	Send(Telemetry)
	Flush()
	Stop()
	IsThrottled() bool
	Close(flush, retry bool, retryTimeout time.Duration) chan bool
}
