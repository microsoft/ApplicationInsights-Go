package appinsights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func telemetryBuffer(items ...Telemetry) TelemetryBufferItems {
	ctx := NewTelemetryContext()
	ctx.iKey = "00000000-0000-0000-0000-000000000000"

	var result TelemetryBufferItems
	for _, item := range items {
		result = append(result, ctx.envelop(item))
	}

	return result
}

func (buffer *TelemetryBufferItems) add(items ...Telemetry) {
	*buffer = append(*buffer, telemetryBuffer(items...)...)
}

func TestJsonSerializerSingle(t *testing.T) {
	mockClock()
	defer resetClock()

	item := NewTraceTelemetry("testing", Verbose)
	nowString := currentClock.Now().Format(time.RFC3339)

	j, err := parsePayload(telemetryBuffer(item).serialize())
	if err != nil {
		t.Errorf("Error parsing payload: %s", err.Error())
	}

	if len(j) != 1 {
		t.Fatal("Unexpected event count")
	}

	j[0].assertPath(t, "name", "Microsoft.ApplicationInsights.Message")
	j[0].assertPath(t, "time", nowString)
	j[0].assertPath(t, "sampleRate", 100)
	j[0].assertPath(t, "data.baseType", "MessageData")
	j[0].assertPath(t, "data.baseData.message", "testing")
	j[0].assertPath(t, "data.baseData.severityLevel", 0)
	j[0].assertPath(t, "data.baseData.ver", 2)
}

func TestJsonSerializerMultiple(t *testing.T) {
	mockClock()
	defer resetClock()

	var buffer TelemetryBufferItems
	now := currentClock.Now()
	nowString := now.Format(time.RFC3339)

	buffer.add(
		NewTraceTelemetry("testing", Error),
		NewEventTelemetry("an-event"),
		NewMetricTelemetry("a-metric", 567),
	)

	req := NewRequestTelemetry("method", "my-url", time.Minute, "204")
	req.Name = "req-name"
	buffer.add(req)

	j, err := parsePayload(buffer.serialize())
	if err != nil {
		t.Errorf("Error parsing payload: %s", err.Error())
	}

	if len(j) != 4 {
		t.Fatal("Unexpected event count")
	}

	// Trace
	j[0].assertPath(t, "name", "Microsoft.ApplicationInsights.Message")
	j[0].assertPath(t, "time", nowString)
	j[0].assertPath(t, "sampleRate", 100.0)
	j[0].assertPath(t, "data.baseType", "MessageData")
	j[0].assertPath(t, "data.baseData.message", "testing")
	j[0].assertPath(t, "data.baseData.severityLevel", 3)
	j[0].assertPath(t, "data.baseData.ver", 2)

	// Event
	j[1].assertPath(t, "name", "Microsoft.ApplicationInsights.Event")
	j[1].assertPath(t, "time", nowString)
	j[1].assertPath(t, "sampleRate", 100.0)
	j[1].assertPath(t, "data.baseType", "EventData")
	j[1].assertPath(t, "data.baseData.name", "an-event")
	j[1].assertPath(t, "data.baseData.ver", 2)

	// Metric
	j[2].assertPath(t, "name", "Microsoft.ApplicationInsights.Metric")
	j[2].assertPath(t, "time", nowString)
	j[2].assertPath(t, "sampleRate", 100.0)
	j[2].assertPath(t, "data.baseType", "MetricData")
	j[2].assertPath(t, "data.baseData.metrics.<len>", 1)
	j[2].assertPath(t, "data.baseData.metrics.[0].value", 567)
	j[2].assertPath(t, "data.baseData.metrics.[0].count", 1)
	j[2].assertPath(t, "data.baseData.metrics.[0].kind", 0)
	j[2].assertPath(t, "data.baseData.ver", 2)

	// Request
	j[3].assertPath(t, "name", "Microsoft.ApplicationInsights.Request")
	j[3].assertPath(t, "time", now.Add(-time.Minute).Format(time.RFC3339)) // Constructor subtracts duration
	j[3].assertPath(t, "sampleRate", 100.0)
	j[3].assertPath(t, "data.baseType", "RequestData")
	j[3].assertPath(t, "data.baseData.name", "req-name")
	j[3].assertPath(t, "data.baseData.duration", "0.00:01:00.0000000")
	j[3].assertPath(t, "data.baseData.responseCode", "204")
	j[3].assertPath(t, "data.baseData.success", true)
	j[3].assertPath(t, "data.baseData.url", "my-url")
	j[3].assertPath(t, "data.baseData.ver", 2)

	if id, err := j[3].getPath("data.baseData.id"); err != nil {
		t.Errorf("Id not present")
	} else if len(id.(string)) == 0 {
		t.Errorf("Empty request id")
	}
}

type jsonMessage map[string]interface{}
type jsonPayload []jsonMessage

func parsePayload(payload []byte) (jsonPayload, error) {
	// json.Decoder can detect line endings for us but I'd like to explicitly find them.
	var result jsonPayload
	for _, item := range bytes.Split(payload, []byte("\n")) {
		if len(item) == 0 {
			continue
		}

		decoder := json.NewDecoder(bytes.NewReader(item))
		msg := make(jsonMessage)
		if err := decoder.Decode(&msg); err == nil {
			result = append(result, msg)
		} else {
			return result, err
		}
	}

	return result, nil
}

func (msg jsonMessage) assertPath(t *testing.T, path string, value interface{}) {
	const tolerance = 0.0000001
	v, err := msg.getPath(path)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if num, ok := value.(int); ok {
		if vnum, ok := v.(float64); ok {
			if math.Abs(float64(num)-vnum) > tolerance {
				t.Errorf("Data was unexpected at %s. Got %g want %d", path, vnum, num)
			}
		} else if vnum, ok := v.(int); ok {
			if vnum != num {
				t.Errorf("Data was unexpected at %s. Got %d want %d", path, vnum, num)
			}
		} else {
			t.Errorf("Expected value at %s to be a number, but was %t", path, v)
		}
	} else if num, ok := value.(float64); ok {
		if vnum, ok := v.(float64); ok {
			if math.Abs(num-vnum) > tolerance {
				t.Errorf("Data was unexpected at %s. Got %g want %g", path, vnum, num)
			}
		} else if vnum, ok := v.(int); ok {
			if math.Abs(num-float64(vnum)) > tolerance {
				t.Errorf("Data was unexpected at %s. Got %d want %g", path, vnum, num)
			}
		} else {
			t.Errorf("Expected value at %s to be a number, but was %t", path, v)
		}
	} else if str, ok := value.(string); ok {
		if vstr, ok := v.(string); ok {
			if str != vstr {
				t.Errorf("Data was unexpected at %s. Got '%s' want '%s'", path, vstr, str)
			}
		} else {
			t.Errorf("Expected value at %s to be a string, but was %t", path, v)
		}
	} else if bl, ok := value.(bool); ok {
		if vbool, ok := v.(bool); ok {
			if bl != vbool {
				t.Errorf("Data was unexpected at %s. Got %q want %q", path, vbool, bl)
			}
		} else {
			t.Errorf("Expected value at %t to be a bool, but was %t", path, v)
		}
	} else {
		t.Errorf("Unsupported type: %v", value)
	}
}

func (msg jsonMessage) getPath(path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	var obj interface{} = msg
	for i, part := range parts {
		if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			// Array
			idxstr := part[1 : len(part)-2]
			idx, _ := strconv.Atoi(idxstr)

			if ar, ok := obj.([]interface{}); ok {
				if idx >= len(ar) {
					return nil, fmt.Errorf("Index out of bounds: %s", strings.Join(parts[0:i+1], "."))
				}

				obj = ar[idx]
			} else {
				return nil, fmt.Errorf("Path %s is not an array", strings.Join(parts[0:i], "."))
			}
		} else if part == "<len>" {
			if ar, ok := obj.([]interface{}); ok {
				return len(ar), nil
			}
		} else {
			// Map
			if dict, ok := obj.(jsonMessage); ok {
				if val, ok := dict[part]; ok {
					obj = val
				} else {
					return nil, fmt.Errorf("Key %s not found in %s", part, strings.Join(parts[0:i], "."))
				}
			} else if dict, ok := obj.(map[string]interface{}); ok {
				if val, ok := dict[part]; ok {
					obj = val
				} else {
					return nil, fmt.Errorf("Key %s not found in %s", part, strings.Join(parts[0:i], "."))
				}
			} else {
				return nil, fmt.Errorf("Path %s is not a map", strings.Join(parts[0:i], "."))
			}
		}
	}

	return obj, nil
}
