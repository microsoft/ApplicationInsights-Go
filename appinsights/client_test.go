package appinsights

import "testing"

func TestClientBurstPerformance(t *testing.T) {
	client := NewTelemetryClient("")
	client.(*telemetryClient).channel.(*InMemoryChannel).transmitter = &nullTransmitter{}

	for i := 0; i < 1000000; i++ {
		client.TrackTrace("A message")
	}
}
