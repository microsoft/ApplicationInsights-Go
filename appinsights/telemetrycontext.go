package appinsights

import (
	"os"
	"runtime"
)

type TelemetryContext struct {
	InstrumentationKey string
	Tags               map[string]string
}

func NewItemTelemetryContext() TelemetryContext {
	return TelemetryContext{
		Tags: make(map[string]string)}
}

func NewClientTelemetryContext() TelemetryContext {
	context := TelemetryContext{
		Tags: make(map[string]string),
	}
	loadDeviceContext(&context)
	loadInternalContext(&context)
	return context
}

func loadDeviceContext(context *TelemetryContext) {
	hostname, err := os.Hostname()
	if err == nil {
		context.Tags[DeviceId] = hostname
		context.Tags[DeviceMachineName] = hostname
		context.Tags[DeviceRoleInstance] = hostname
	}
	context.Tags[DeviceOS] = runtime.GOOS
}

func loadInternalContext(context *TelemetryContext) {
	context.Tags[InternalSdkVersion] = "go:" + version
}
