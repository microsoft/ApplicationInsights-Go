package appinsights

type TelemetryContext struct {
	iKey string
	Tags map[string]string
}

func NewTelemetryContext *TelemetryContext {
	return &TelemetryContext{
		Tags: make(map[string]string),
	}
}

func (context *telemetryContext) getStringTag(key string) string {
	if result, ok := context.tags[key]; ok {
		return result
	}

	return ""
}

func (context *telemetryContext) setStringTag(key, value string) {
	if value != "" {
		context.tags[key] = value
	} else {
		delete(context.tags, key)
	}
}

func (context *telemetryContext) getBoolTag(key string) bool {
	if result, ok := context.tags[key]; ok {
		if value, err := strconv.ParseBool(result); err == nil {
			return value
		}
	}

	return false
}

func (context *telemetryContext) setBoolTag(key string, value bool) {
	if value {
		context.tags[key] = "true"
	} else {
		delete(context.tags, key)
	}
}
