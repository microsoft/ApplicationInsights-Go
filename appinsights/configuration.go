package appinsights

type TelemetryConfiguration struct {
	InstrumentationKey string
	EndpointUrl        string
}

func NewTelemetryConfiguration(instrumentationKey string) *TelemetryConfiguration {
	return &TelemetryConfiguration{
		InstrumentationKey: instrumentationKey,
		EndpointUrl:        "https://dc.services.visualstudio.com/v2/track",
	}
}
