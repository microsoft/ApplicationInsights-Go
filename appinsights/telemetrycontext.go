package appinsights

import (
	"strconv"
	"time"

	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/satori/go.uuid"
)

// The telemetry context type stores context keys that will be applied to
// submitted telemetry.  This includes, e.g. information about the system
// and application sending the telemetry as well as information used for
// correlation with other events.  Each TelemetryClient contains a
// TelemetryContext that will set values on every outgoing item if the key
// is not overridden inside the telemetry item's TelemetryContext.
type TelemetryContext struct {
	// Instrumentation key
	iKey string

	// Collection of tag data to attach to the telemetry item.
	Tags map[string]string

	// Common properties to add to each telemetry item.  This only has
	// an effect from the TelemetryClient's context instance.  This will
	// be nil on telemetry items.
	CommonProperties map[string]string
}

// Creates a new, empty TelemetryContext
func NewTelemetryContext() *TelemetryContext {
	return &TelemetryContext{
		Tags: make(map[string]string),
	}
}

// Gets the instrumentation key associated with this TelemetryContext.  This
// will be an empty string on telemetry items' context instances.
func (context *TelemetryContext) InstrumentationKey() string {
	return context.iKey
}

// Wraps a telemetry item in an envelope with the information found in this
// context.
func (context *TelemetryContext) envelop(item Telemetry) *contracts.Envelope {
	// Apply common properties
	if props := item.GetProperties(); props != nil && context.CommonProperties != nil {
		for k, v := range context.CommonProperties {
			if _, ok := props[k]; !ok {
				props[k] = v
			}
		}
	}

	tdata := item.TelemetryData()
	data := contracts.NewData()
	data.BaseType = tdata.BaseType()
	data.BaseData = tdata

	envelope := contracts.NewEnvelope()
	envelope.Name = tdata.EnvelopeName()
	envelope.Data = data
	envelope.IKey = context.iKey

	timestamp := item.Time()
	if timestamp.IsZero() {
		timestamp = currentClock.Now()
	}

	envelope.Time = timestamp.UTC().Format(time.RFC3339)

	if itemContext := item.TelemetryContext(); itemContext != nil && itemContext != context {
		envelope.Tags = itemContext.Tags

		// Copy in default tag values.
		for tagkey, tagval := range context.Tags {
			if _, ok := itemContext.Tags[tagkey]; !ok {
				envelope.Tags[tagkey] = tagval
			}
		}
	} else {
		// Create new tags object
		envelope.Tags = make(map[string]string)
		for k, v := range context.Tags {
			envelope.Tags[k] = v
		}
	}

	// Create operation ID if it does not exist
	if _, ok := envelope.Tags[contracts.OperationId]; !ok {
		envelope.Tags[contracts.OperationId] = uuid.NewV4().String()
	}

	// Sanitize.
	for _, warn := range tdata.Sanitize() {
		diagnosticsWriter.Printf("Telemetry data warning: %s", warn)
	}
	for _, warn := range contracts.SanitizeTags(envelope.Tags) {
		diagnosticsWriter.Printf("Telemetry tag warning: %s", warn)
	}

	return envelope
}

func (context *TelemetryContext) getStringTag(key string) string {
	if result, ok := context.Tags[key]; ok {
		return result
	}

	return ""
}

func (context *TelemetryContext) setStringTag(key, value string) {
	if value != "" {
		context.Tags[key] = value
	} else {
		delete(context.Tags, key)
	}
}

func (context *TelemetryContext) getBoolTag(key string) bool {
	if result, ok := context.Tags[key]; ok {
		if value, err := strconv.ParseBool(result); err == nil {
			return value
		}
	}

	return false
}

func (context *TelemetryContext) setBoolTag(key string, value bool) {
	if value {
		context.Tags[key] = "true"
	} else {
		delete(context.Tags, key)
	}
}
