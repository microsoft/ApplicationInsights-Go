package appinsights

import (
	"os"
	"runtime"
	"time"
)

type TelemetryConfiguration struct {
	InstrumentationKey string
	EndpointUrl        string
	MaxBatchSize       int
	MaxBatchInterval   time.Duration
}

func NewTelemetryConfiguration(instrumentationKey string) *TelemetryConfiguration {
	return &TelemetryConfiguration{
		InstrumentationKey: instrumentationKey,
		EndpointUrl:        "https://dc.services.visualstudio.com/v2/track",
		MaxBatchSize:       1024,
		MaxBatchInterval:   time.Duration(10) * time.Second,
	}
}

func (config *TelemetryConfiguration) setupContext(context *telemetryContext) {
	context.iKey = config.InstrumentationKey
	context.Internal().SetSdkVersion("go:" + version)
	context.Device().SetOsVersion(runtime.GOOS)

	if hostname, err := os.Hostname(); err == nil {
		context.Device().SetId(hostname)
		context.Cloud().SetRoleInstance(hostname)
	}
}
