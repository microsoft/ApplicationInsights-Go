package appinsights

import "time"

type TelemetryClient struct {
	TelemetryConfiguration TelemetryConfiguration
	channel                TelemetryChannel
	Context                TelemetryContext
	IsEnabled              bool
}

func NewTelemetryClient(iKey string) TelemetryClient {
	config := NewTelemetryConfiguration(iKey)
	return TelemetryClient{
		TelemetryConfiguration: config,
		channel:                NewInMemoryChannel(config.EndpointUrl),
		Context:                NewClientTelemetryContext(),
		IsEnabled:              true,
	}
}

func (tc TelemetryClient) InstrumentationKey() string {
	return tc.TelemetryConfiguration.InstrumentationKey
}

func (tc *TelemetryClient) Track(item Telemetry) {
	if !tc.IsEnabled {
		return
	}

	iKey := tc.Context.InstrumentationKey
	if len(iKey) == 0 {
		iKey = tc.TelemetryConfiguration.InstrumentationKey
	}

	item.Context.InstrumentationKey = iKey

	for tagkey, tagval := range tc.Context.Tags {
		if item.Context.Tags[tagkey] == "" {
			item.Context.Tags[tagkey] = tagval
		}
	}

	tc.channel.Send(item)
}

func (tc *TelemetryClient) TrackEvent(name string) {
	tc.Track(NewEventTelemetry(name))
}

func (tc *TelemetryClient) TrackMetric(name string, value float32) {
	tc.Track(NewMetricTelemetry(name, value))
}

func (tc *TelemetryClient) TrackTrace(message string) {
	tc.Track(NewTraceTelemetry(message, Information))
}

func (tc *TelemetryClient) TrackRequest(name string, timestamp time.Time, duration time.Duration, responseCode string, success bool) {
	tc.Track(NewRequestTelemetry(name, timestamp, duration, responseCode, success))
}
