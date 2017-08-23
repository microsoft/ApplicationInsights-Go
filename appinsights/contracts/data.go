package contracts

// NOTE: This file was automatically generated.

// Data struct to contain both B and C sections.
type Data struct {
	Base

	// Container for data item (B section).
	BaseData interface{} `json:"baseData"`
}

// Creates a new Data instance with default values set by the schema.
func NewData() *Data {
	return &Data{}
}
