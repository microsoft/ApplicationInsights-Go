package appinsights

import "bytes"
import "encoding/json"
import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "time"

type InMemoryChannel interface {
	EndpointAddress() string
	Send(Telemetry)
}

type inMemoryChannel struct {
	endpointAddress string
}

func NewInMemoryChannel(endpointAddress string) InMemoryChannel {
	return &inMemoryChannel{
		endpointAddress: endpointAddress,
	}
}

func (channel *inMemoryChannel) EndpointAddress() string {
	return channel.endpointAddress
}

func (channel *inMemoryChannel) Send(item Telemetry) {    
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
		log.Fatal(err)
		return
	}

	buf := bytes.NewReader(jsonBytes)

	req, err := http.NewRequest("POST", channel.EndpointAddress(), buf)
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("Application Insights Telemetry: %s", string(jsonBytes))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println()
	fmt.Printf("Response: %d", resp.StatusCode)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf(" - %s", body)
	fmt.Println()
}
