package appinsights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type TelemetryBufferItems []Telemetry

type InMemoryChannel interface {
	EndpointAddress() string
	Send(Telemetry)
	Flush()
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
		for _ = range ticker.C {
			//log.Trace("Transmit tick at ", t)
			channel.Flush()
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
	defer channel.bufferWg.Done()
	channel.buffer = append(channel.buffer, item)
}

func (channel *inMemoryChannel) Requeue(items TelemetryBufferItems) {
	for _, item := range items {
		channel.Send(item)
	}
}

func (channel *inMemoryChannel) swapBuffer() TelemetryBufferItems {
	channel.bufferWg.Add(1)
	defer channel.bufferWg.Done()
	buffer := channel.buffer
	channel.buffer = make(TelemetryBufferItems, 0)
	return buffer
}

func (channel *inMemoryChannel) Flush() {
	if len(channel.buffer) == 0 {
		//log.Trace("Not transmitting due to empty buffer.")
		return
	}

	buffer := channel.swapBuffer()

	transmission := fmt.Sprintf("\n----------- Transmitting %d items ---------\n\n", len(buffer))

	// TODO: Return the actual buffer here instead of buffer -> string -> buffer
	reqBody := buffer.serialize()
	reqBuf := bytes.NewBufferString(reqBody)

	req, err := http.NewRequest("POST", channel.endpointAddress, reqBuf)
	if err != nil {
		log.Printf("complete requeuing, due to \"%s\"", err)
		channel.Requeue(buffer)
		log.Print("requeuing done.")
		return
	}

	req.Header.Set("Content-Type", "application/x-json-stream")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	transmission += fmt.Sprintf(reqBody)

	start := time.Now()
	client := http.DefaultClient
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil || resp == nil {
		if err != nil {
			log.Printf("complete requeuing, due to \"%s\"", err)
		} else if resp == nil {
			log.Print("complete requeuing, due to missing response")
		}
		channel.Requeue(buffer)
		log.Print("requeuing done.")
		return
	}

	transmission += fmt.Sprintf("\nSent in %s\n", duration)
	transmission += fmt.Sprintf("Response: %d", resp.StatusCode)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("requeuing, due to \"%s\"", err)
		channel.Requeue(buffer)
		log.Print("requeuing done.")
		return
	}

	var report struct {
		Errors []struct {
			Index      int    `json:"index"`
			Message    string `json:"message"`
			StatusCode int    `json:"statusCode"`
		} `json:"errors"`
	}

	err = json.Unmarshal(body, &report)
	if err != nil {
		log.Printf("requeuing, due to \"%s\"", err)
		channel.Requeue(buffer)
		log.Print("requeuing done.")
		return
	}

	for _, reportedError := range report.Errors {
		log.Printf(
			"requeuing index %d, due to status code %d and message \"%s\"",
			reportedError.Index,
			reportedError.StatusCode,
			reportedError.Message)
		channel.Send(buffer[reportedError.Index])
	}

	transmission += fmt.Sprintf(" - %s\n", body)
	transmission += fmt.Sprintf("\n-----------------------------------------\n")

	diagWriter.Write(transmission)
}
