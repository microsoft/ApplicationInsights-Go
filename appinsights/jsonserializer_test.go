package appinsights

import (
	"fmt"
	"testing"
	"time"
)

func TestJsonSerializerSingle(t *testing.T) {
	item := NewTraceTelemetry("testing", Verbose)
	now := time.Now()
	item.Timestamp = now

	want := fmt.Sprintf(`{"ver":1,"name":"Microsoft.ApplicationInsights.Message","time":"%s","sampleRate":100,"seq":"","iKey":"","data":{"baseType":"MessageData","baseData":{"ver":2,"message":"testing","severityLevel":0}}}`, now.Format(time.RFC3339))
	want += "\n"

	items := TelemetryBufferItems{item}
	buf := items.serialize()

	if string(buf) != want {
		t.Errorf("serialize() returned %q, want %q", string(buf), want)
	}
}

func TestJsonSerializerMultiple(t *testing.T) {
	buffer := make(TelemetryBufferItems, 0)
	now := time.Now()
	nowString := now.Format(time.RFC3339)

	for i := 0; i < 3; i++ {
		item := NewTraceTelemetry("testing", Verbose)
		item.Timestamp = now
		buffer = append(buffer, item)
	}

	want := fmt.Sprintf(`{"ver":1,"name":"Microsoft.ApplicationInsights.Message","time":"%s","sampleRate":100,"seq":"","iKey":"","data":{"baseType":"MessageData","baseData":{"ver":2,"message":"testing","severityLevel":0}}}`+"\n"+`{"ver":1,"name":"Microsoft.ApplicationInsights.Message","time":"%s","sampleRate":100,"seq":"","iKey":"","data":{"baseType":"MessageData","baseData":{"ver":2,"message":"testing","severityLevel":0}}}`+"\n"+`{"ver":1,"name":"Microsoft.ApplicationInsights.Message","time":"%s","sampleRate":100,"seq":"","iKey":"","data":{"baseType":"MessageData","baseData":{"ver":2,"message":"testing","severityLevel":0}}}`+"\n",
		nowString,
		nowString,
		nowString)
	result := buffer.serialize()
	if string(result) != want {
		t.Errorf("serialize() returned '%s', want '%s'", result, want)
	}
}
