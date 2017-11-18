package appinsights

import (
	"testing"
	"time"
)

func BenchmarkClientBurstPerformance(b *testing.B) {
	client := NewTelemetryClient("")
	client.(*telemetryClient).channel.(*InMemoryChannel).transmitter = &nullTransmitter{}

	for i := 0; i < b.N; i++ {
		client.TrackTrace("A message")
	}

	<-client.Channel().Close(time.Minute)
}

func TestDefaultTags(t *testing.T) {
	client := NewTelemetryClient("")
	client.Context().Tags["test"] = "OK"
	telem := NewTraceTelemetry("Hello world.", Verbose)
	envelope := client.Context().envelop(telem)
	if envelope.Tags["test"] != "OK" {
		t.Error("Default client tags did not propagate to telemetry")
	}
}
