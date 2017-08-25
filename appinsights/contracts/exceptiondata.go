package contracts

// NOTE: This file was automatically generated.

// An instance of Exception represents a handled or unhandled exception that
// occurred during execution of the monitored application.
type ExceptionData struct {
	Domain

	// Schema version
	Ver int `json:"ver"`

	// Exception chain - list of inner exceptions.
	Exceptions []*ExceptionDetails `json:"exceptions"`

	// Severity level. Mostly used to indicate exception severity level when it is
	// reported by logging library.
	SeverityLevel SeverityLevel `json:"severityLevel"`

	// Identifier of where the exception was thrown in code. Used for exceptions
	// grouping. Typically a combination of exception type and a function from the
	// call stack.
	ProblemId string `json:"problemId"`

	// Collection of custom properties.
	Properties map[string]string `json:"properties"`

	// Collection of custom measurements.
	Measurements map[string]float64 `json:"measurements"`
}

// Returns the name used when this is embedded within an Envelope container.
func (data *ExceptionData) EnvelopeName() string {
	return "Microsoft.ApplicationInsights.Exception"
}

// Returns the base type when placed within a Data object container.
func (data *ExceptionData) BaseType() string {
	return "ExceptionData"
}

// Creates a new ExceptionData instance with default values set by the schema.
func NewExceptionData() *ExceptionData {
	return &ExceptionData{
		Ver:          2,
		Properties:   make(map[string]string),
		Measurements: make(map[string]float64),
	}
}
