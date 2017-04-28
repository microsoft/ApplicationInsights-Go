package appinsights

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	submit_retries = []time.Duration{time.Duration(10 * time.Second), time.Duration(30 * time.Second), time.Duration(60 * time.Second)}
)

type TelemetryBufferItems []Telemetry

type InMemoryChannel struct {
	endpointAddress string
	isDeveloperMode bool
	collectChan     chan Telemetry
	flushChan       chan bool
	batchSize       int
	batchInterval   time.Duration
}

func NewInMemoryChannel(config *TelemetryConfiguration) *InMemoryChannel {
	channel := &InMemoryChannel{
		endpointAddress: config.EndpointUrl,
		collectChan:     make(chan Telemetry),
		flushChan:       make(chan bool),
		batchSize:       config.MaxBatchSize,
		batchInterval:   config.MaxBatchInterval,
	}

	go channel.acceptLoop()

	return channel
}

func (channel *InMemoryChannel) EndpointAddress() string {
	return channel.endpointAddress
}

func (channel *InMemoryChannel) Send(item Telemetry) {
	if item != nil {
		channel.collectChan <- item
	}
}

func (channel *InMemoryChannel) Flush() {
	channel.flushChan <- true
}

func (channel *InMemoryChannel) Stop() {
	close(channel.collectChan)
}

func (channel *InMemoryChannel) acceptLoop() {
	buffer := make(TelemetryBufferItems, 16)

	for {
		if len(buffer) > 16 {
			// Start out with the size of the previous buffer
			buffer = make(TelemetryBufferItems, len(buffer))
		} else if len(buffer) > 0 {
			// Start out with at least 16 slots
			buffer = make(TelemetryBufferItems, 16)
		}

		// Wait for an event
		select {
		case event := <-channel.collectChan:
			if event == nil {
				// Channel closed, quit.
				close(channel.flushChan)
				return
			}

			buffer = append(buffer, event)

		case _ = <-channel.flushChan:
			// The buffer is empty.
			break
		}

		if len(buffer) == 0 {
			continue
		}

		// Delay until timeout passes
		timer := time.NewTimer(channel.batchInterval)
	waitLoop:
		for {
			select {
			case event := <-channel.collectChan:
				if event == nil {
					// Channel closed, flush and exit.
					break waitLoop
				}

				buffer = append(buffer, event)
				if len(buffer) >= channel.batchSize {
					break waitLoop
				}

			case _ = <-channel.flushChan:
				break waitLoop

			case _ = <-timer.C:
				// Timeout expired
				break waitLoop
			}
		}

		if !timer.Stop() {
			<-timer.C
		}

		if len(buffer) > 0 {
			go channel.transmitRetry(buffer)
		}
	}
}

func (channel *InMemoryChannel) transmitRetry(items TelemetryBufferItems) {
	var diagnostics bytes.Buffer
	body := items.serialize()

	for _, wait := range submit_retries {
		diagnostics.Reset()
		err := channel.transmit(len(items), body, &diagnostics)
		diagnosticsWriter.Write(diagnostics.String())

		if err == nil {
			return
		}

		time.Sleep(wait)
	}

	// One final try
	diagnostics.Reset()
	err := channel.transmit(len(items), body, &diagnostics)
	diagnosticsWriter.Write(diagnostics.String())
	if err != nil {
		diagnosticsWriter.Write("Gave up transmitting payload; exhausted retries")
	}
}

func (channel *InMemoryChannel) transmit(count int, body []byte, diag *bytes.Buffer) error {
	fmt.Fprintf(diag, "\n----------- Transmitting %d items ---------\n\n", count)

	start := time.Now()

	req, err := http.NewRequest("POST", channel.endpointAddress, bytes.NewReader(body))
	if err != nil {
		// Requeue
		fmt.Fprintf(diag, "Error from NewRequest: %s", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/x-json-stream")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	diag.Write(body)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(diag, "\nError from client.Do: %s", err.Error())
		return err
	}

	duration := time.Since(start)

	fmt.Fprintf(diag, "\nSent in %s\n", duration)
	fmt.Fprintf(diag, "Response: %d", resp.StatusCode)

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(diag, "Error reading response: %s", err.Error())
		return err
	}

	fmt.Fprintf(diag, " - %s\n", respBody)
	fmt.Fprintf(diag, "\n-----------------------------------------\n")

	return nil
}
