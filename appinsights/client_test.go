package appinsights

import "testing"

func TestClientBurstPerformance(t *testing.T) {
	telemetryClient := NewTelemetryClient("")
	for i := 0; i < 1000000; i++ {
		telemetryClient.TrackTrace("A message")
	}
}
