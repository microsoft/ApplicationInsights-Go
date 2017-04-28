package appinsights

type TelemetryChannel interface {
	EndpointAddress() string
	Send(Telemetry)
	Flush()
	Stop()
}
