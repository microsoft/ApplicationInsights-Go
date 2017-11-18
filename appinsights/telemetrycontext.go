package appinsights

import (
	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
	"strconv"
	"time"
)

type TelemetryContext struct {
	iKey string
	Tags map[string]string
}

func NewTelemetryContext() *TelemetryContext {
	return &TelemetryContext{
		Tags: make(map[string]string),
	}
}

func (context *TelemetryContext) InstrumentationKey() string {
	return context.iKey
}

func (context *TelemetryContext) envelop(item Telemetry) *contracts.Envelope {
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

	envelope.Time = timestamp.Format(time.RFC3339)

	if itemContext := item.TelemetryContext(); context != nil {
		envelope.Tags = itemContext.Tags

		// Copy in default tag values.
		if context != itemContext {
			for tagkey, tagval := range context.Tags {
				if _, ok := itemContext.Tags[tagkey]; !ok {
					envelope.Tags[tagkey] = tagval
				}
			}
		}
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
