package appinsights

import "testing"

func TestClientBurstPerformance(t *testing.T) {
	config := NewTelemetryConfiguration("")
	config.EndpointUrl = ""
	telemetryClient := NewTelemetryClientFromConfig(config)

	for i := 0; i < 1000000; i++ {
		telemetryClient.TrackTrace("A message")
	}
}
