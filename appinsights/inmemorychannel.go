package appinsights

import "bytes"
import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "time"
import "sync"

type TelemetryBufferItems []Telemetry

type InMemoryChannel interface {
	EndpointAddress() string
	Send(Telemetry)
}

type inMemoryChannel struct {
	endpointAddress string
	isDeveloperMode bool
	buffer          TelemetryBufferItems
	bufferWg        sync.WaitGroup
	ticker          *time.Ticker
}

var diagWriter = getDiagnosticsMessageWriter()

func NewInMemoryChannel(endpointAddress string) InMemoryChannel {
	buffer := make(TelemetryBufferItems, 0)
	var bufferWg sync.WaitGroup
	ticker := time.NewTicker(time.Second * 10)
	channel := &inMemoryChannel{
		endpointAddress: endpointAddress,
		buffer:          buffer,
		bufferWg:        bufferWg,
		ticker:          ticker,
	}

	go func() {
		for t := range ticker.C {
			//log.Trace("Transmit tick at ", t)
			channel.transmit(t)
		}
	}()

	return channel
}

func (channel *inMemoryChannel) EndpointAddress() string {
	return channel.endpointAddress
}

func (channel *inMemoryChannel) Send(item Telemetry) {
	// TODO: Use a fixed buffer size and don't require sync.
	channel.bufferWg.Add(1)
	channel.buffer = append(channel.buffer, item)
	channel.bufferWg.Done()
}

func (channel *inMemoryChannel) swapBuffer() TelemetryBufferItems {
	channel.bufferWg.Add(1)
	buffer := channel.buffer
	channel.buffer = make(TelemetryBufferItems, 0)
	channel.bufferWg.Done()
	return buffer
}

func (channel *inMemoryChannel) transmit(t time.Time) {
	if len(channel.buffer) == 0 {
		//log.Trace("Not transmitting due to empty buffer.")
		return
	}

	buffer := channel.swapBuffer()

	transmission := fmt.Sprintf("\n----------- Transmitting %d items ---------\n\n", len(buffer))

	start := time.Now()

	// TODO: Return the actual buffer here instead of buffer -> string -> buffer
	reqBody := buffer.serialize()
	reqBuf := bytes.NewBufferString(reqBody)

	req, err := http.NewRequest("POST", channel.endpointAddress, reqBuf)
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-json-stream")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	transmission += fmt.Sprintf(reqBody)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	duration := time.Since(start)

	transmission += fmt.Sprintf("\nSent in %s\n", duration)
	transmission += fmt.Sprintf("Response: %d", resp.StatusCode)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	transmission += fmt.Sprintf(" - %s\n", body)
	transmission += fmt.Sprintf("\n-----------------------------------------\n")

	diagWriter.Write(transmission)
}
