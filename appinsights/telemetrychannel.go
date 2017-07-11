package appinsights

import "time"

type TelemetryChannel interface {
	EndpointAddress() string
	Send(Telemetry)
	Flush()
	Stop()
	Close(flush bool, retryTimeout time.Duration) chan bool
	CloseSync(flush bool, retryTimeout time.Duration)
}
