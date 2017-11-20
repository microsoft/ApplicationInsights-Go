package contracts

// NOTE: This file was automatically generated.

// The abstract common base of all domains.
type Domain struct {
}

func (data *Domain) Sanitize() []string {
	var warnings []string

	return warnings
}

// Creates a new Domain instance with default values set by the schema.
func NewDomain() *Domain {
	return &Domain{}
}
