package contracts

// NOTE: This file was automatically generated.

// Instances of Event represent structured event records that can be grouped
// and searched by their properties. Event data item also creates a metric of
// event count by name.
type EventData struct {
	Domain

	// Schema version
	Ver int `json:"ver"`

	// Event name. Keep it low cardinality to allow proper grouping and useful
	// metrics.
	Name string `json:"name"`

	// Collection of custom properties.
	Properties map[string]string `json:"properties"`

	// Collection of custom measurements.
	Measurements map[string]float64 `json:"measurements"`
}

// Creates a new EventData instance with default values set by the schema.
func NewEventData() *EventData {
	return &EventData{
		Ver:          2,
		Properties:   make(map[string]string),
		Measurements: make(map[string]float64),
	}
}

// Returns the name used when this is embedded within an Envelope container.
func (data *EventData) EnvelopeName() string {
	return "Microsoft.ApplicationInsights.Event"
}

// Returns the base type when placed within a Data object container.
func (data *EventData) BaseType() string {
	return "EventData"
}
