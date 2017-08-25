package appinsights

import (
	"fmt"
	"testing"
	"time"
)

func TestJsonSerializerSingle(t *testing.T) {
	item := NewTraceTelemetry("testing", Verbose)
	now := time.Now()
	item.timestamp = now

	want := fmt.Sprintf(`{"name":"Microsoft.ApplicationInsights.Message","time":"%s","iKey":"","tags":{},"data":{"baseType":"MessageData","baseData":{"ver":2,"properties":null,"message":"testing","severityLevel":0}}}`, now.Format(time.RFC3339))
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
		item.timestamp = now
		buffer = append(buffer, item)
	}

	want := fmt.Sprintf(`{"name":"Microsoft.ApplicationInsights.Message","time":"%s","iKey":"","tags":{},"data":{"baseType":"MessageData","baseData":{"ver":2,"properties":null,"message":"testing","severityLevel":0}}}`+"\n"+`{"name":"Microsoft.ApplicationInsights.Message","time":"%s","iKey":"","tags":{},"data":{"baseType":"MessageData","baseData":{"ver":2,"properties":null,"message":"testing","severityLevel":0}}}`+"\n"+`{"name":"Microsoft.ApplicationInsights.Message","time":"%s","iKey":"","tags":{},"data":{"baseType":"MessageData","baseData":{"ver":2,"properties":null,"message":"testing","severityLevel":0}}}`+"\n",
		nowString,
		nowString,
		nowString)
	result := buffer.serialize()
	if string(result) != want {
		t.Errorf("serialize() returned '%s', want '%s'", result, want)
	}
}
