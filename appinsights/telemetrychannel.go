package appinsights

type TelemetryChannel interface {
	Send(Telemetry)
}
