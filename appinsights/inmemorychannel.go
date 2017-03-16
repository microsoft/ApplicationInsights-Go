package appinsights

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"sync"
)

type TelemetryBufferItems []Telemetry

type InMemoryChannel interface {
	EndpointAddress() string
	Send(Telemetry)
}

type inMemoryChannel struct {
	endpointAddress string
	isDeveloperMode bool
	buffer          TelemetryBufferItems
	bufferLock      sync.Mutex
	ticker          *time.Ticker
}

var diagWriter = getDiagnosticsMessageWriter()

func NewInMemoryChannel(endpointAddress string) InMemoryChannel {
	buffer := make(TelemetryBufferItems, 0)
	ticker := time.NewTicker(time.Second * 10)
	channel := &inMemoryChannel{
		endpointAddress: endpointAddress,
		buffer:          buffer,
		ticker:          ticker,
	}

	go func() {
		for _ = range ticker.C {
			channel.transmit()
		}
	}()

	return channel
}

func (channel *inMemoryChannel) EndpointAddress() string {
	return channel.endpointAddress
}

func (channel *inMemoryChannel) Send(item Telemetry) {
	// TODO: Use a fixed buffer size and don't require sync.
	channel.bufferLock.Lock()
	defer channel.bufferLock.Unlock()
	channel.buffer = append(channel.buffer, item)
}

func (channel *inMemoryChannel) sendMany(items TelemetryBufferItems) {
	channel.bufferLock.Lock()
	defer channel.bufferLock.Unlock()

	channel.buffer = append(channel.buffer, items...)
}

func (channel *inMemoryChannel) swapBuffer() TelemetryBufferItems {
	channel.bufferLock.Lock()
	defer channel.bufferLock.Unlock()

	buffer := channel.buffer
	channel.buffer = make(TelemetryBufferItems, 0)
	return buffer
}

func (channel *inMemoryChannel) transmit() error {
	if len(channel.buffer) == 0 {
		//log.Trace("Not transmitting due to empty buffer.")
		return nil
	}

	buffer := channel.swapBuffer()

	var dbg bytes.Buffer
	fmt.Fprintf(&dbg, "\n----------- Transmitting %d items ---------\n\n", len(buffer))

	start := time.Now()

	// TODO: Return the actual buffer here instead of buffer -> string -> buffer
	reqBody := buffer.serialize()
	reqBuf := bytes.NewBufferString(reqBody)

	req, err := http.NewRequest("POST", channel.endpointAddress, reqBuf)
	if err != nil {
		// Requeue
		fmt.Fprintf(&dbg, "Error from NewRequest: %s", err.Error())
		diagWriter.Write(dbg.String())
		channel.sendMany(buffer)
		return err
	}

	req.Header.Set("Content-Type", "application/x-json-stream")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	dbg.WriteString(reqBody)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(&dbg, "\nError from client.Do: %s", err.Error())
		diagWriter.Write(dbg.String())
		channel.sendMany(buffer)
		return err
	}

	duration := time.Since(start)

	fmt.Fprintf(&dbg, "\nSent in %s\n", duration)
	fmt.Fprintf(&dbg, "Response: %d", resp.StatusCode)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(&dbg, "Error reading response: %s", err.Error())
		diagWriter.Write(dbg.String())
		channel.sendMany(buffer)
		return err
	}

	fmt.Fprintf(&dbg, " - %s\n", body)
	fmt.Fprintf(&dbg, "\n-----------------------------------------\n")
	diagWriter.Write(dbg.String())

	return nil
}
