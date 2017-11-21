package appinsights

import (
	"time"

	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
)

type TelemetryClient interface {
	Context() *TelemetryContext
	InstrumentationKey() string
	Channel() TelemetryChannel
	IsEnabled() bool
	SetIsEnabled(bool)
	Track(Telemetry)
	TrackEvent(string)
	TrackMetric(string, float64)
	TrackTrace(string, contracts.SeverityLevel)
	TrackRequest(string, string, time.Duration, string)
	TrackRemoteDependency(name, dependencyType, target string, success bool)
	TrackAvailability(name string, duration time.Duration, success bool)
	TrackPageView(name, url string)
	TrackException(err interface{})
}

type telemetryClient struct {
	channel                TelemetryChannel
	context                *TelemetryContext
	isEnabled              bool
}

func NewTelemetryClient(iKey string) TelemetryClient {
	return NewTelemetryClientFromConfig(NewTelemetryConfiguration(iKey))
}

func NewTelemetryClientFromConfig(config *TelemetryConfiguration) TelemetryClient {
	channel := NewInMemoryChannel(config)
	context := NewTelemetryContext()

	config.setupContext(context)

	return &telemetryClient{
		channel:                channel,
		context:                context,
		isEnabled:              true,
	}
}

func (tc *telemetryClient) Context() *TelemetryContext {
	return tc.context
}

func (tc *telemetryClient) Channel() TelemetryChannel {
	return tc.channel
}

func (tc *telemetryClient) InstrumentationKey() string {
	return tc.context.InstrumentationKey()
}

func (tc *telemetryClient) IsEnabled() bool {
	return tc.isEnabled
}

func (tc *telemetryClient) SetIsEnabled(isEnabled bool) {
	tc.isEnabled = isEnabled
}

func (tc *telemetryClient) Track(item Telemetry) {
	if tc.isEnabled {
		tc.channel.Send(tc.context.envelop(item))
	}
}

func (tc *telemetryClient) TrackEvent(name string) {
	tc.Track(NewEventTelemetry(name))
}

func (tc *telemetryClient) TrackMetric(name string, value float64) {
	tc.Track(NewMetricTelemetry(name, value))
}

func (tc *telemetryClient) TrackTrace(message string, severity contracts.SeverityLevel) {
	tc.Track(NewTraceTelemetry(message, severity))
}

func (tc *telemetryClient) TrackRequest(method, url string, duration time.Duration, responseCode string) {
	tc.Track(NewRequestTelemetry(method, url, duration, responseCode))
}

func (tc *telemetryClient) TrackRemoteDependency(name, dependencyType, target string, success bool) {
	tc.Track(NewRemoteDependencyTelemetry(name, dependencyType, target, success))
}

func (tc *telemetryClient) TrackAvailability(name string, duration time.Duration, success bool) {
	tc.Track(NewAvailabilityTelemetry(name, duration, success))
}

func (tc *telemetryClient) TrackPageView(name, url string) {
	tc.Track(NewPageViewTelemetry(name, url))
}

func (tc *telemetryClient) TrackException(err interface{}) {
	tc.Track(NewExceptionTelemetry(err))
}
