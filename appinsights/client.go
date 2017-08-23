package appinsights

import (
	"time"

	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
)

type TelemetryClient interface {
	Context() TelemetryContext
	InstrumentationKey() string
	Channel() TelemetryChannel
	IsEnabled() bool
	SetIsEnabled(bool)
	Track(Telemetry)
	TrackEvent(string)
	TrackMetric(string, float64)
	TrackTrace(string)
	TrackRequest(string, string, string, time.Time, time.Duration, string, bool)
}

type telemetryClient struct {
	TelemetryConfiguration *TelemetryConfiguration
	channel                TelemetryChannel
	context                TelemetryContext
	isEnabled              bool
}

func NewTelemetryClient(iKey string) TelemetryClient {
	return NewTelemetryClientFromConfig(NewTelemetryConfiguration(iKey))
}

func NewTelemetryClientFromConfig(config *TelemetryConfiguration) TelemetryClient {
	channel := NewInMemoryChannel(config)
	context := NewTelemetryContext()

	config.setupContext(context.(*telemetryContext))

	return &telemetryClient{
		TelemetryConfiguration: config,
		channel:                channel,
		context:                context,
		isEnabled:              true,
	}
}

func (tc *telemetryClient) Context() TelemetryContext {
	return tc.context
}

func (tc *telemetryClient) Channel() TelemetryChannel {
	return tc.channel
}

func (tc *telemetryClient) InstrumentationKey() string {
	return tc.TelemetryConfiguration.InstrumentationKey
}

func (tc *telemetryClient) IsEnabled() bool {
	return tc.isEnabled
}

func (tc *telemetryClient) SetIsEnabled(isEnabled bool) {
	tc.isEnabled = isEnabled
}

func (tc *telemetryClient) Track(item Telemetry) {
	if tc.isEnabled {
		iKey := tc.context.InstrumentationKey()

		itemContext := item.Context().(*telemetryContext)
		itemContext.iKey = iKey

		if clientContext, ok := tc.context.(*telemetryContext); ok {
			for tagkey, tagval := range clientContext.tags {
				if _, ok := itemContext.tags[tagkey]; !ok {
					itemContext.tags[tagkey] = tagval
				}
			}
		}

		tc.channel.Send(item)
	}
}

func (tc *telemetryClient) TrackEvent(name string) {
	tc.Track(NewEventTelemetry(name))
}

func (tc *telemetryClient) TrackMetric(name string, value float64) {
	tc.Track(NewMetricTelemetry(name, value))
}

func (tc *telemetryClient) TrackTrace(message string) {
	tc.Track(NewTraceTelemetry(message, contracts.Information))
}

func (tc *telemetryClient) TrackRequest(name, method, url string, timestamp time.Time, duration time.Duration, responseCode string, success bool) {
	tc.Track(NewRequestTelemetry(name, method, url, timestamp, duration, responseCode, success))
}
