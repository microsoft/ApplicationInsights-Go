package contracts

// NOTE: This file was automatically generated.

// Exception details of the exception in a chain.
type ExceptionDetails struct {

	// In case exception is nested (outer exception contains inner one), the id
	// and outerId properties are used to represent the nesting.
	Id int `json:"id"`

	// The value of outerId is a reference to an element in ExceptionDetails that
	// represents the outer exception
	OuterId int `json:"outerId"`

	// Exception type name.
	TypeName string `json:"typeName"`

	// Exception message.
	Message string `json:"message"`

	// Indicates if full exception stack is provided in the exception. The stack
	// may be trimmed, such as in the case of a StackOverflow exception.
	HasFullStack bool `json:"hasFullStack"`

	// Text describing the stack. Either stack or parsedStack should have a value.
	Stack string `json:"stack"`

	// List of stack frames. Either stack or parsedStack should have a value.
	ParsedStack []*StackFrame `json:"parsedStack,omitempty"`
}

// Creates a new ExceptionDetails instance with default values set by the schema.
func NewExceptionDetails() *ExceptionDetails {
	return &ExceptionDetails{
		HasFullStack: true,
	}
}
