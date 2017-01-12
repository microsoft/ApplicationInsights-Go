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
	envelope := Envelope{
		Name: "Microsoft.ApplicationInsights." + item.TypeName,
		Time: item.Timestamp.Format(time.RFC3339),
		IKey: item.Context.InstrumentationKey,
		Data: Data{
			BaseType: item.TypeName + "Data",
			BaseData: item.Data},
		Tags: item.Context.Tags}

	jsonBytes, err := json.Marshal(envelope)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(jsonBytes)
}
