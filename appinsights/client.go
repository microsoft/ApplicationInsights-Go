package appinsights

import "time"

type TelemetryClient struct {
	TelemetryConfiguration TelemetryConfiguration
	Channel                TelemetryChannel
	Context                TelemetryContext
	IsEnabled              bool
}

func NewTelemetryClient(iKey string) TelemetryClient {
	config := NewTelemetryConfiguration(iKey)
	return TelemetryClient{
		TelemetryConfiguration: config,
		Channel:                NewInMemoryChannel(config.EndpointUrl),
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

	tc.Channel.Send(item)
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

func (tc *TelemetryClient) TrackRequest(
	id string,
	name string,
	timestamp time.Time,
	duration time.Duration,
	httpMethod string,
	url string,
	responseCode string,
	success bool) {

	tc.Track(NewRequestTelemetry(
		id,
		name,
		timestamp,
		duration,
		httpMethod,
		url,
		responseCode,
		success))
}

func (tc *TelemetryClient) TrackRemoteDependency(
	id string,
	name string,
	resultCode int,
	commandName string,
	kind DataPointType,
	duration time.Duration,
	count int,
	min float32,
	max float32,
	stdDev float32,
	theType string,
	dependencyKind DependencyKind,
	success bool,
	async bool,
	dependencySource DependencySourceType,
	properties map[string]string,
	alterTelementryContext func(*TelemetryContext)) {

	tc.Track(NewRemoteDependencyData(
		id,
		name,
		resultCode,
		commandName,
		kind,
		duration,
		count,
		min,
		max,
		stdDev,
		theType,
		dependencyKind,
		success,
		async,
		dependencySource,
		properties,
		alterTelementryContext))
}
