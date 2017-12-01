package appinsights

import (
	"strings"
	"testing"

	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
)

func TestDefaultTags(t *testing.T) {
	context := NewTelemetryContext()
	context.Tags["test"] = "OK"
	context.Tags["no-write"] = "Fail"

	telem := NewTraceTelemetry("Hello world.", Verbose)
	telem.Context.Tags["no-write"] = "OK"

	envelope := context.envelop(telem)

	if envelope.Tags["test"] != "OK" {
		t.Error("Default client tags did not propagate to telemetry")
	}

	if envelope.Tags["no-write"] != "OK" {
		t.Error("Default client tag overwrote telemetry item tag")
	}
}

func TestCommonProperties(t *testing.T) {
	context := NewTelemetryContext()
	context.CommonProperties = map[string]string{
		"test":     "OK",
		"no-write": "Fail",
	}

	telem := NewTraceTelemetry("Hello world.", Verbose)
	telem.Properties["no-write"] = "OK"

	envelope := context.envelop(telem)
	data := envelope.Data.(*contracts.Data).BaseData.(*contracts.MessageData)

	if data.Properties["test"] != "OK" {
		t.Error("Common properties did not propagate to telemetry")
	}

	if data.Properties["no-write"] != "OK" {
		t.Error("Common properties overwrote telemetry properties")
	}
}

func TestTagHelpers(t *testing.T) {
	context := NewTelemetryContext()
	if context.getStringTag("Nonexistent") != "" {
		t.Error("Successfully fetched nonexistent tag")
	}

	context.setStringTag("my_tag", "foo")
	if v, ok := context.Tags["my_tag"]; !ok || v != "foo" {
		t.Error("setStringTag had no effect")
	}

	if context.getStringTag("my_tag") != "foo" {
		t.Error("setStringTag/getStringTag had unexpected result")
	}

	context.setStringTag("my_tag", "")
	if _, ok := context.Tags["my_tag"]; ok {
		t.Error("setStringTag did not delete tag")
	}

	if context.getBoolTag("Nonexistent") != false {
		t.Error("Getting a nonexistent bool tag did not default to false")
	}

	context.setBoolTag("my_bool", true)
	if v, ok := context.Tags["my_bool"]; !ok || v != "true" {
		t.Error("Setting a bool tag should set it to 'true'")
	}

	if context.getBoolTag("my_bool") != true {
		t.Error("Getting a bool should return the correct value")
	}

	context.setBoolTag("my_bool", false)
	if _, ok := context.Tags["my_bool"]; ok {
		t.Error("Setting a tag to false should remove it from the dict")
	}

	if context.getBoolTag("my_bool") != false {
		t.Error("Getting a bool should return the correct value")
	}
}

func TestSanitize(t *testing.T) {
	name := strings.Repeat("Z", 1024)
	val := strings.Repeat("Y", 10240)

	ev := NewEventTelemetry(name)
	ev.Properties[name] = val
	ev.Measurements[name] = 55.0

	ctx := NewTelemetryContext()
	ctx.Session().SetId(name)

	// We'll be looking for messages with these values:
	found := map[string]int{
		"EventData.Name exceeded":        0,
		"EventData.Properties has value": 0,
		"EventData.Properties has key":   0,
		"EventData.Measurements has key": 0,
		"ai.session.id exceeded":         0,
	}

	// Set up listener for the warnings.
	NewDiagnosticsMessageListener(func(msg string) error {
		for k, _ := range found {
			if strings.Contains(msg, k) {
				found[k] = found[k] + 1
				break
			}
		}

		return nil
	})

	defer resetDiagnosticsListeners()

	// This may break due to hardcoded limits... Check contracts.
	envelope := ctx.envelop(ev)

	// Make sure all the warnings were found in the output
	for k, v := range found {
		if v != 1 {
			t.Errorf("Did not find a warning containing \"%s\"", k)
		}
	}

	// Check the format of the stuff we found in the envelope
	if v, ok := envelope.Tags[contracts.SessionId]; !ok || v != name[:64] {
		t.Error("Session ID tag was not truncated")
	}

	evdata := envelope.Data.(*contracts.Data).BaseData.(*contracts.EventData)
	if evdata.Name != name[:512] {
		t.Error("Event name was not truncated")
	}

	if v, ok := evdata.Properties[name[:150]]; !ok || v != val[:8192] {
		t.Error("Event property name/value was not truncated")
	}

	if v, ok := evdata.Measurements[name[:150]]; !ok || v != 55.0 {
		t.Error("Event measurement name was not truncated")
	}
}
