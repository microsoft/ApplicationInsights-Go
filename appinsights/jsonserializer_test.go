package appinsights

import "fmt"
import "testing"
import "time"

func TestJsonSerializerSingle(t *testing.T) {

	item := NewTraceTelemetry("testing", nil, Verbose)
	now := time.Now()
	item.Timestamp = now

	want := fmt.Sprintf(`{"name":"Microsoft.ApplicationInsights.Message","time":"%s","iKey":"","tags":{},"data":{"baseType":"MessageData","baseData":{"ver":2,"properties":null,"message":"testing","severityLevel":0}}}`, now.Format(time.RFC3339))
	result := serialize(item)
	if result != want {
		t.Errorf("serialize() returned '%s', want '%s'", result, want)
	}
}

func TestJsonSerializerMultiple(t *testing.T) {

	buffer := make(TelemetryBufferItems, 0)
	now := time.Now()
	nowString := now.Format(time.RFC3339)

	for i := 0; i < 3; i++ {
		item := NewTraceTelemetry("testing", nil, Verbose)
		item.Timestamp = now
		buffer = append(buffer, item)
	}

	want := fmt.Sprintf(`{"name":"Microsoft.ApplicationInsights.Message","time":"%s","iKey":"","tags":{},"data":{"baseType":"MessageData","baseData":{"ver":2,"properties":null,"message":"testing","severityLevel":0}}}`+"\n"+`{"name":"Microsoft.ApplicationInsights.Message","time":"%s","iKey":"","tags":{},"data":{"baseType":"MessageData","baseData":{"ver":2,"properties":null,"message":"testing","severityLevel":0}}}`+"\n"+`{"name":"Microsoft.ApplicationInsights.Message","time":"%s","iKey":"","tags":{},"data":{"baseType":"MessageData","baseData":{"ver":2,"properties":null,"message":"testing","severityLevel":0}}}`+"\n",
		nowString,
		nowString,
		nowString)
	result := buffer.serialize()
	if result != want {
		t.Errorf("serialize() returned '%s', want '%s'", result, want)
	}
}
