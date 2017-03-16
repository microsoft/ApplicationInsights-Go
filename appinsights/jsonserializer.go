package appinsights

import "bytes"
import "encoding/json"
import "log"
import "time"

func (items TelemetryBufferItems) serialize() string {
	var result bytes.Buffer

	for i := range items {
		item := items[i]
		result.WriteString(serialize(item))
		result.WriteString("\n")
	}

	return result.String()
}

func serialize(item Telemetry) string {
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

	envelope.Tags = context.(*telemetryContext).tags

	jsonBytes, err := json.Marshal(envelope)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(jsonBytes)
}
