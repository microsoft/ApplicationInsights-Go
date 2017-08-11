package appinsights

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

const ten_seconds = time.Duration(10) * time.Second

type testTransmitter struct {
	requests  chan *testTransmission
	responses chan *transmissionResult
}

func (transmitter *testTransmitter) Transmit(payload []byte, items TelemetryBufferItems) (*transmissionResult, error) {
	itemsCopy := make(TelemetryBufferItems, len(items))
	copy(itemsCopy, items)

	transmitter.requests <- &testTransmission{
		payload:   string(payload),
		items:     itemsCopy,
		timestamp: currentClock.Now(),
	}

	return <-transmitter.responses, nil
}

func (transmitter *testTransmitter) Close() {
	close(transmitter.requests)
	close(transmitter.responses)
}

func (transmitter *testTransmitter) prepResponse(statusCodes ...int) {
	for _, code := range statusCodes {
		transmitter.responses <- &transmissionResult{statusCode: code}
	}
}

func (transmitter *testTransmitter) waitForRequest(t *testing.T) *testTransmission {
	select {
	case req := <-transmitter.requests:
		return req
	case <-time.After(time.Duration(500) * time.Millisecond):
		t.Fatal("Timed out waiting for request to be sent")
		return nil /* Not reached */
	}
}

func (transmitter *testTransmitter) assertNoRequest(t *testing.T) {
	select {
	case <-transmitter.requests:
		t.Fatal("Expected no request")
	case <-time.After(time.Duration(10) * time.Millisecond):
		return
	}
}

type testTransmission struct {
	timestamp time.Time
	payload   string
	items     TelemetryBufferItems
}

func newTestChannelServer(config ...*TelemetryConfiguration) (TelemetryClient, *testTransmitter) {
	transmitter := &testTransmitter{
		requests:  make(chan *testTransmission, 16),
		responses: make(chan *transmissionResult, 16),
	}

	var client TelemetryClient
	if len(config) > 0 {
		client = NewTelemetryClientFromConfig(config[0])
	} else {
		config := NewTelemetryConfiguration("")
		config.MaxBatchInterval = ten_seconds // assumed by every test.
		client = NewTelemetryClientFromConfig(config)
	}

	client.(*telemetryClient).channel.(*InMemoryChannel).transmitter = transmitter

	return client, transmitter
}

func assertTimeApprox(t *testing.T, x, y time.Time) {
	const delta = (time.Duration(100) * time.Millisecond)
	if (x.Before(y) && y.Sub(x) > delta) || (y.Before(x) && x.Sub(y) > delta) {
		t.Errorf("Time isn't a close match: %v vs %v", x, y)
	}
}

func slowInc(seconds int) {
	const delay = time.Millisecond * time.Duration(5)

	// Sleeps in tests are evil, but with all the async nonsense going
	// on, no callbacks, and minimal control of the clock, I'm not
	// really sure I have another choice.

	time.Sleep(delay)
	for i := 0; i < seconds; i++ {
		fakeClock.Increment(time.Second)
		time.Sleep(delay)
	}
}

func waitForClose(t *testing.T, ch chan bool) bool {
	select {
	case res := <-ch:
		return res
	case <-time.After(time.Duration(100) * time.Second):
		t.Fatal("Close signal not received in 100ms")
		return false /* not reached */
	}
}

func TestSimpleSubmit(t *testing.T) {
	mockClock()
	defer resetClock()
	client, transmitter := newTestChannelServer()
	defer transmitter.Close()
	defer client.Channel().Close(false, false, 0)

	client.TrackTrace("~msg~")
	tm := currentClock.Now()
	transmitter.prepResponse(200)

	slowInc(11)
	req := transmitter.waitForRequest(t)

	assertTimeApprox(t, req.timestamp, tm.Add(ten_seconds))

	if !strings.Contains(string(req.payload), "~msg~") {
		t.Errorf("Payload does not contain message")
	}
}

func TestMultipleSubmit(t *testing.T) {
	mockClock()
	defer resetClock()
	client, transmitter := newTestChannelServer()
	defer transmitter.Close()
	defer client.Channel().Close(false, false, 0)

	transmitter.prepResponse(200, 200)

	start := currentClock.Now()

	for i := 0; i < 16; i++ {
		client.TrackTrace(fmt.Sprintf("~msg-%x~", i))
		slowInc(1)
	}

	slowInc(10)

	req1 := transmitter.waitForRequest(t)
	assertTimeApprox(t, req1.timestamp, start.Add(ten_seconds))

	for i := 0; i < 10; i++ {
		if !strings.Contains(req1.payload, fmt.Sprintf("~msg-%x~", i)) {
			t.Errorf("Payload does not contain expected item: %x", i)
		}
	}

	req2 := transmitter.waitForRequest(t)
	assertTimeApprox(t, req2.timestamp, start.Add(ten_seconds+ten_seconds))

	for i := 10; i < 16; i++ {
		if !strings.Contains(req2.payload, fmt.Sprintf("~msg-%x~", i)) {
			t.Errorf("Payload does not contain expected item: %x", i)
		}
	}
}

func TestFlush(t *testing.T) {
	mockClock()
	defer resetClock()
	client, transmitter := newTestChannelServer()
	defer transmitter.Close()
	defer client.Channel().Close(false, false, 0)

	transmitter.prepResponse(200, 200)

	// Empty flush should do nothing
	client.Channel().Flush()

	tm := currentClock.Now()
	client.TrackTrace("~msg~")
	client.Channel().Flush()

	req1 := transmitter.waitForRequest(t)
	assertTimeApprox(t, req1.timestamp, tm)
	if !strings.Contains(req1.payload, "~msg~") {
		t.Error("Unexpected payload")
	}

	// Next one goes back to normal
	client.TrackTrace("~next~")
	slowInc(11)

	req2 := transmitter.waitForRequest(t)
	assertTimeApprox(t, req2.timestamp, tm.Add(ten_seconds))
	if !strings.Contains(req2.payload, "~next~") {
		t.Error("Unexpected payload")
	}
}

func TestCloseNoFlush(t *testing.T) {
	mockClock()
	defer resetClock()
	client, transmitter := newTestChannelServer()
	defer transmitter.Close()

	transmitter.prepResponse(200)

	client.TrackTrace("Not sent")
	client.Channel().Close(false, false, 0)
	slowInc(20)
	transmitter.assertNoRequest(t)
}

func TestCloseFlush(t *testing.T) {
	mockClock()
	defer resetClock()
	client, transmitter := newTestChannelServer()
	defer transmitter.Close()

	transmitter.prepResponse(200)

	client.TrackTrace("~flushed~")
	client.Channel().Close(true, false, 0)

	req := transmitter.waitForRequest(t)
	if !strings.Contains(req.payload, "~flushed~") {
		t.Error("Unexpected payload")
	}
}

func TestCloseFlushRetry(t *testing.T) {
	mockClock()
	defer resetClock()
	client, transmitter := newTestChannelServer()
	defer transmitter.Close()

	transmitter.prepResponse(500, 200)

	client.TrackTrace("~flushed~")
	tm := currentClock.Now()
	ch := client.Channel().Close(true, true, time.Minute)

	slowInc(30)

	waitForClose(t, ch)

	req1 := transmitter.waitForRequest(t)
	if !strings.Contains(req1.payload, "~flushed~") {
		t.Error("Unexpected payload")
	}

	assertTimeApprox(t, req1.timestamp, tm)

	req2 := transmitter.waitForRequest(t)
	if !strings.Contains(req2.payload, "~flushed~") {
		t.Error("Unexpected payload")
	}

	assertTimeApprox(t, req2.timestamp, tm.Add(submit_retries[0]))
}

func TestCloseWithOngoingRetry(t *testing.T) {
	mockClock()
	defer resetClock()
	client, transmitter := newTestChannelServer()
	defer transmitter.Close()

	transmitter.prepResponse(408, 200, 200)

	// This message should get stuck, retried
	client.TrackTrace("~msg-1~")
	slowInc(11)

	// This message will get flushed immediately
	client.TrackTrace("~msg-2~")
	ch := client.Channel().Close(true, true, time.Minute)

	// Then, let's wait for the first message to go out...
	slowInc(30)

	waitForClose(t, ch)

	// Check.
	req1 := transmitter.waitForRequest(t)
	if !strings.Contains(req1.payload, "~msg-1~") {
		t.Error("First message unexpected payload")
	}

	req2 := transmitter.waitForRequest(t)
	if !strings.Contains(req2.payload, "~msg-2~") {
		t.Error("Second message unexpected payload")
	}

	req3 := transmitter.waitForRequest(t)
	if !strings.Contains(req3.payload, "~msg-1~") {
		t.Error("Third message unexpected payload")
	}
}
