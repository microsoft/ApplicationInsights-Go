package appinsights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func (items TelemetryBufferItems) serialize() []byte {
	var result bytes.Buffer

	for _, item := range items {
		end := result.Len()
		if err := serialize(item, &result); err != nil {
			diagnosticsWriter.Write(fmt.Sprintf("Telemetry item failed to serialize: %s", err.Error()))
			result.Truncate(end)
		}
	}

	return result.Bytes()
}

func serialize(item Telemetry, writer io.Writer) error {
	data := &data{
		BaseType: item.baseTypeName() + "Data",
		BaseData: item.baseData(),
	}

	context := item.Context()

	envelope := &envelope{
		Name: "Microsoft.ApplicationInsights." + item.baseTypeName(),
		Time: item.Timestamp().Format(time.RFC3339),
		IKey: context.InstrumentationKey(),
		Data: data,
	}

	if tcontext, ok := context.(*telemetryContext); ok {
		envelope.Tags = tcontext.tags
	}

	return json.NewEncoder(writer).Encode(envelope)
}
