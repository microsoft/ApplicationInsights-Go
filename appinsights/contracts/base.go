package contracts

// NOTE: This file was automatically generated.

// Data struct to contain only C section with custom fields.
type Base struct {

	// Name of item (B section) if any. If telemetry data is derived straight from
	// this, this should be null.
	BaseType string `json:"baseType"`
}

// Creates a new Base instance with default values set by the schema.
func NewBase() *Base {
	return &Base{}
}
