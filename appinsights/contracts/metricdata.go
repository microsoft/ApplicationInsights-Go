package contracts

// NOTE: This file was automatically generated.

// An instance of the Metric item is a list of measurements (single data
// points) and/or aggregations.
type MetricData struct {
	Domain

	// Schema version
	Ver int `json:"ver"`

	// List of metrics. Only one metric in the list is currently supported by
	// Application Insights storage. If multiple data points were sent only the
	// first one will be used.
	Metrics []*DataPoint `json:"metrics"`

	// Collection of custom properties.
	Properties map[string]string `json:"properties,omitempty"`
}

// Returns the name used when this is embedded within an Envelope container.
func (data *MetricData) EnvelopeName() string {
	return "Microsoft.ApplicationInsights.Metric"
}

// Returns the base type when placed within a Data object container.
func (data *MetricData) BaseType() string {
	return "MetricData"
}

// Creates a new MetricData instance with default values set by the schema.
func NewMetricData() *MetricData {
	return &MetricData{
		Ver:        2,
		Properties: make(map[string]string),
	}
}
