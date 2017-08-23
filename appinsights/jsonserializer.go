package appinsights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jjjordanmsft/ApplicationInsights-Go/appinsights/contracts"
)

func (items TelemetryBufferItems) serialize() []byte {
	var result bytes.Buffer
	encoder := json.NewEncoder(&result)

	for _, item := range items {
		end := result.Len()
		if err := encoder.Encode(prepare(item)); err != nil {
			diagnosticsWriter.Write(fmt.Sprintf("Telemetry item failed to serialize: %s", err.Error()))
			result.Truncate(end)
		}
	}

	return result.Bytes()
}

func prepare(item Telemetry) *contracts.Envelope {
	context := item.Context()
	tdata := item.TelemetryData()

	data := contracts.NewData()
	data.BaseType = tdata.BaseType()
	data.BaseData = tdata

	envelope := contracts.NewEnvelope()
	envelope.Name = tdata.EnvelopeName()
	envelope.Time = item.Time().Format(time.RFC3339)
	envelope.IKey = context.InstrumentationKey()
	envelope.Data = data

	if tcontext, ok := context.(*telemetryContext); ok {
		envelope.Tags = tcontext.tags
	}

	return envelope
}
