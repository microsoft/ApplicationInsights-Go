package appinsights

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	submit_telemetry_every = time.Duration(10 * time.Second)
	submit_retries = []time.Duration{time.Duration(10 * time.Second), time.Duration(30 * time.Second), time.Duration(60 * time.Second)}
)

type TelemetryBufferItems []Telemetry

type InMemoryChannel interface {
	EndpointAddress() string
	Send(Telemetry)
	Flush()
	Stop()
}

type inMemoryChannel struct {
	endpointAddress string
	isDeveloperMode bool
	collectChan     chan Telemetry
	flushChan       chan bool
}

func NewInMemoryChannel(endpointAddress string) InMemoryChannel {
	channel := &inMemoryChannel{
		endpointAddress: endpointAddress,
		collectChan:     make(chan Telemetry),
		flushChan:       make(chan bool),
	}

	go channel.acceptLoop()

	return channel
}

func (channel *inMemoryChannel) EndpointAddress() string {
	return channel.endpointAddress
}

func (channel *inMemoryChannel) Send(item Telemetry) {
	if item != nil {
		channel.collectChan <- item
	}
}

func (channel *inMemoryChannel) Flush() {
	channel.flushChan <- true
}

func (channel *inMemoryChannel) Stop() {
	close(channel.collectChan)
}

func (channel *inMemoryChannel) acceptLoop() {
	var buffer TelemetryBufferItems

	for {
		// Wait for an event
		select {
		case event := <- channel.collectChan:
			if event == nil {
				// Channel closed, quit.
				close(channel.flushChan)
				return
			}

			buffer = append(buffer, event)

		case _ = <- channel.flushChan:
			// The buffer is empty.
			break
		}

		if len(buffer) == 0 {
			continue
		}

		// Delay until timeout passes
		timer := time.After(submit_telemetry_every)
waitLoop:	for {
			select {
			case event := <- channel.collectChan:
				if event == nil {
					// Channel closed, flush and exit.
					break waitLoop
				}

				buffer = append(buffer, event)

			case _ = <- channel.flushChan:
				break waitLoop

			case _ = <- timer:
				// Timeout expired
				break waitLoop
			}
		}

		reqBody := buffer.serialize()
		count := len(buffer)
		buffer = buffer[:0]

		if len(reqBody) > 0 {
			go channel.transmitRetry(count, reqBody)
		}
	}
}

func (channel *inMemoryChannel) transmitRetry(count int, reqBody []byte) {
	for _, wait := range submit_retries {
		if err := channel.transmit(count, reqBody); err == nil {
			return
		}
		
		time.Sleep(wait)
	}
	
	diagnosticsWriter.Write("Gave up transmitting payload; exhausted retries")
}

func (channel *inMemoryChannel) transmit(count int, reqBody []byte) error {
	var dbg bytes.Buffer
	fmt.Fprintf(&dbg, "\n----------- Transmitting %d items ---------\n\n", count)

	start := time.Now()

	req, err := http.NewRequest("POST", channel.endpointAddress, bytes.NewReader(reqBody))
	if err != nil {
		// Requeue
		fmt.Fprintf(&dbg, "Error from NewRequest: %s", err.Error())
		diagnosticsWriter.Write(dbg.String())
		return err
	}

	req.Header.Set("Content-Type", "application/x-json-stream")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	dbg.Write(reqBody)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(&dbg, "\nError from client.Do: %s", err.Error())
		diagnosticsWriter.Write(dbg.String())
		return err
	}

	duration := time.Since(start)

	fmt.Fprintf(&dbg, "\nSent in %s\n", duration)
	fmt.Fprintf(&dbg, "Response: %d", resp.StatusCode)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(&dbg, "Error reading response: %s", err.Error())
		diagnosticsWriter.Write(dbg.String())
		return err
	}

	fmt.Fprintf(&dbg, " - %s\n", body)
	fmt.Fprintf(&dbg, "\n-----------------------------------------\n")
	diagnosticsWriter.Write(dbg.String())

	return nil
}
