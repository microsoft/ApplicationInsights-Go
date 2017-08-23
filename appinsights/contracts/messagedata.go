package contracts

// NOTE: This file was automatically generated.

// Instances of Message represent printf-like trace statements that are
// text-searched. Log4Net, NLog and other text-based log file entries are
// translated into intances of this type. The message does not have
// measurements.
type MessageData struct {
	Domain

	// Schema version
	Ver int `json:"ver"`

	// Trace message
	Message string `json:"message"`

	// Trace severity level.
	SeverityLevel SeverityLevel `json:"severityLevel"`

	// Collection of custom properties.
	Properties map[string]string `json:"properties"`
}

// Creates a new MessageData instance with default values set by the schema.
func NewMessageData() *MessageData {
	return &MessageData{
		Ver:        2,
		Properties: make(map[string]string),
	}
}

// Returns the name used when this is embedded within an Envelope container.
func (data *MessageData) EnvelopeName() string {
	return "Microsoft.ApplicationInsights.Message"
}

// Returns the base type when placed within a Data object container.
func (data *MessageData) BaseType() string {
	return "MessageData"
}
