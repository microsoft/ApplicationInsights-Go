package appinsights

import "os"
import "runtime"

type TelemetryContext interface {
	InstrumentationKey() string
	loadDeviceContext()
}

type telemetryContext struct {
	iKey string
	tags map[string]string
}

func NewTelemetryContext() TelemetryContext {
	context := &telemetryContext{
		tags: make(map[string]string),
	}
	context.loadDeviceContext()
	context.loadInternalContext()
	return context
}

func (context *telemetryContext) InstrumentationKey() string {
	return context.iKey
}

func (context *telemetryContext) loadDeviceContext() {
	hostname, err := os.Hostname()
	if err == nil {
		context.tags[DeviceId] = hostname
		context.tags[DeviceMachineName] = hostname
		context.tags[DeviceRoleInstance] = hostname
	}
	context.tags[DeviceOS] = runtime.GOOS
}

func (context *telemetryContext) loadInternalContext() {
	context.tags[InternalSdkVersion] = "go:" + version
}
