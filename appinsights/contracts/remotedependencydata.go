package contracts

// NOTE: This file was automatically generated.

// An instance of Remote Dependency represents an interaction of the monitored
// component with a remote component/service like SQL or an HTTP endpoint.
type RemoteDependencyData struct {
	Domain

	// Schema version
	Ver int `json:"ver"`

	// Name of the command initiated with this dependency call. Low cardinality
	// value. Examples are stored procedure name and URL path template.
	Name string `json:"name"`

	// Identifier of a dependency call instance. Used for correlation with the
	// request telemetry item corresponding to this dependency call.
	Id string `json:"id"`

	// Result code of a dependency call. Examples are SQL error code and HTTP
	// status code.
	ResultCode string `json:"resultCode"`

	// Request duration in format: DD.HH:MM:SS.MMMMMM. Must be less than 1000
	// days.
	Duration string `json:"duration"`

	// Indication of successfull or unsuccessfull call.
	Success bool `json:"success"`

	// Command initiated by this dependency call. Examples are SQL statement and
	// HTTP URL's with all query parameters.
	Data string `json:"data"`

	// Target site of a dependency call. Examples are server name, host address.
	Target string `json:"target"`

	// Dependency type name. Very low cardinality value for logical grouping of
	// dependencies and interpretation of other fields like commandName and
	// resultCode. Examples are SQL, Azure table, and HTTP.
	Type string `json:"type"`

	// Collection of custom properties.
	Properties map[string]string `json:"properties"`

	// Collection of custom measurements.
	Measurements map[string]float64 `json:"measurements"`
}

// Creates a new RemoteDependencyData instance with default values set by the schema.
func NewRemoteDependencyData() *RemoteDependencyData {
	return &RemoteDependencyData{
		Ver:          2,
		Success:      true,
		Properties:   make(map[string]string),
		Measurements: make(map[string]float64),
	}
}

// Returns the name used when this is embedded within an Envelope container.
func (data *RemoteDependencyData) EnvelopeName() string {
	return "Microsoft.ApplicationInsights.RemoteDependency"
}

// Returns the base type when placed within a Data object container.
func (data *RemoteDependencyData) BaseType() string {
	return "RemoteDependencyData"
}
