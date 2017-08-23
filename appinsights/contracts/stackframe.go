package contracts

// NOTE: This file was automatically generated.

// Stack frame information.
type StackFrame struct {

	// Level in the call stack. For the long stacks SDK may not report every
	// function in a call stack.
	Level int `json:"level"`

	// Method name.
	Method string `json:"method"`

	// Name of the assembly (dll, jar, etc.) containing this function.
	Assembly string `json:"assembly"`

	// File name or URL of the method implementation.
	FileName string `json:"fileName"`

	// Line number of the code implementation.
	Line int `json:"line"`
}

// Creates a new StackFrame instance with default values set by the schema.
func NewStackFrame() *StackFrame {
	return &StackFrame{}
}
