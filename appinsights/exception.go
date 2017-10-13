package appinsights

import (
	"fmt"
	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
	"reflect"
	"runtime"
	"strings"
)

type ExceptionTelemetry struct {
	BaseTelemetry
	Error         interface{}
	Frames        []*contracts.StackFrame
	SeverityLevel contracts.SeverityLevel
}

func NewExceptionTelemetry(err interface{}) *ExceptionTelemetry {
	return &ExceptionTelemetry{
		Error:         err,
		Frames:        GetCallstack(2),
		SeverityLevel: Error,
		BaseTelemetry: BaseTelemetry{
			Timestamp:    currentClock.Now(),
			Context:      NewTelemetryContext(),
			Properties:   make(map[string]string),
			Measurements: make(map[string]float64),
		},
	}
}

func (telem *ExceptionTelemetry) TelemetryData() TelemetryData {
	details := contracts.NewExceptionDetails()
	details.HasFullStack = len(telem.Frames) > 0
	details.ParsedStack = telem.Frames

	if err, ok := telem.Error.(error); ok {
		details.Message = err.Error()
		details.TypeName = reflect.TypeOf(telem.Error).String()
	} else if str, ok := telem.Error.(string); ok {
		details.Message = str
		details.TypeName = "string"
	} else if stringer, ok := telem.Error.(fmt.Stringer); ok {
		details.Message = stringer.String()
		details.TypeName = reflect.TypeOf(telem.Error).String()
	} else {
		details.Message = "<unknown>"
		details.TypeName = "<unknown>"
	}

	data := contracts.NewExceptionData()
	data.SeverityLevel = telem.SeverityLevel
	data.Exceptions = []*contracts.ExceptionDetails{details}
	data.Properties = telem.Properties
	data.Measurements = telem.Measurements

	return data
}

func GetCallstack(skip int) []*contracts.StackFrame {
	var stackFrames []*contracts.StackFrame

	stack := make([]uintptr, 64)
	depth := runtime.Callers(skip+1, stack)
	if depth == 0 {
		return stackFrames
	}

	frames := runtime.CallersFrames(stack[:depth])
	level := 0
	for {
		frame, more := frames.Next()

		stackFrame := &contracts.StackFrame{
			Level:    level,
			FileName: frame.File,
			Line:     frame.Line,
		}

		if frame.Function != "" {
			/* Default */
			stackFrame.Method = frame.Function

			/* Break up function into assembly/function */
			lastSlash := strings.LastIndexByte(frame.Function, '/')
			if lastSlash >= 0 {
				firstDot := strings.IndexByte(frame.Function[lastSlash:], '.')
				if firstDot >= 0 {
					stackFrame.Assembly = frame.Function[:lastSlash+firstDot]
					stackFrame.Method = frame.Function[lastSlash+firstDot+1:]
				}
			}
		}

		stackFrames = append(stackFrames, stackFrame)

		level++
		if !more {
			break
		}
	}

	return stackFrames
}
