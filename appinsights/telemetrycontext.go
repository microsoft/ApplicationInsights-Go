package appinsights

import "strconv"

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
