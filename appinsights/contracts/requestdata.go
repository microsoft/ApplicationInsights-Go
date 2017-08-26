package contracts

// NOTE: This file was automatically generated.

// An instance of Request represents completion of an external request to the
// application to do work and contains a summary of that request execution and
// the results.
type RequestData struct {
	Domain

	// Schema version
	Ver int `json:"ver"`

	// Identifier of a request call instance. Used for correlation between request
	// and other telemetry items.
	Id string `json:"id"`

	// Source of the request. Examples are the instrumentation key of the caller
	// or the ip address of the caller.
	Source string `json:"source"`

	// Name of the request. Represents code path taken to process request. Low
	// cardinality value to allow better grouping of requests. For HTTP requests
	// it represents the HTTP method and URL path template like 'GET
	// /values/{id}'.
	Name string `json:"name"`

	// Request duration in format: DD.HH:MM:SS.MMMMMM. Must be less than 1000
	// days.
	Duration string `json:"duration"`

	// Result of a request execution. HTTP status code for HTTP requests.
	ResponseCode string `json:"responseCode"`

	// Indication of successfull or unsuccessfull call.
	Success bool `json:"success"`

	// Request URL with all query string parameters.
	Url string `json:"url"`

	// Collection of custom properties.
	Properties map[string]string `json:"properties,omitempty"`

	// Collection of custom measurements.
	Measurements map[string]float64 `json:"measurements,omitempty"`
}

// Returns the name used when this is embedded within an Envelope container.
func (data *RequestData) EnvelopeName() string {
	return "Microsoft.ApplicationInsights.Request"
}

// Returns the base type when placed within a Data object container.
func (data *RequestData) BaseType() string {
	return "RequestData"
}

// Creates a new RequestData instance with default values set by the schema.
func NewRequestData() *RequestData {
	return &RequestData{
		Ver:          2,
		Properties:   make(map[string]string),
		Measurements: make(map[string]float64),
	}
}
