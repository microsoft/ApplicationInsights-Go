package appinsights

import (
	"sync"
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
	controlChan     chan *inMemoryChannelControl
	batchSize       int
	batchInterval   time.Duration
	waitgroup       sync.WaitGroup
}

type inMemoryChannelControl struct {
	flush    bool
	stop     bool
	timeout  time.Duration
	callback chan bool
}

func NewInMemoryChannel(config *TelemetryConfiguration) *InMemoryChannel {
	channel := &InMemoryChannel{
		endpointAddress: config.EndpointUrl,
		collectChan:     make(chan Telemetry),
		controlChan:     make(chan *inMemoryChannelControl),
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
	channel.controlChan <- &inMemoryChannelControl{
		flush: true,
	}
}

func (channel *InMemoryChannel) Stop() {
	channel.controlChan <- &inMemoryChannelControl{
		stop: true,
	}
}

func (channel *InMemoryChannel) Close(flush bool, timeout time.Duration) chan bool {
	callback := make(chan bool)

	channel.controlChan <- &inMemoryChannelControl{
		stop:     true,
		flush:    flush,
		timeout:  timeout,
		callback: callback,
	}

	return callback
}

func (channel *InMemoryChannel) CloseSync(flush bool, timeout time.Duration) {
	callback := channel.Close(flush, timeout)
	<-callback
}

func (channel *InMemoryChannel) acceptLoop() {
	buffer := make(TelemetryBufferItems, 0, 16)
	stopping := false

mainLoop:
	for !stopping {
		if len(buffer) > 16 {
			// Start out with the size of the previous buffer
			buffer = make(TelemetryBufferItems, 0, cap(buffer))
		} else if len(buffer) > 0 {
			// Start out with at least 16 slots
			buffer = make(TelemetryBufferItems, 0, 16)
		}

		// Wait for an event
		select {
		case event := <-channel.collectChan:
			if event == nil {
				// Channel closed?  Not intercepted by Send()?
				panic("Received nil event")
			}

			buffer = append(buffer, event)

		case ctl := <-channel.controlChan:
			// The buffer is empty, so there would be no point in flushing
			if ctl.stop {
				stopping = true
			}
			if ctl.callback != nil {
				ctl.callback <- true
				close(ctl.callback)
			}
		}

		if len(buffer) == 0 {
			continue
		}

		// Things that are used by the sender if we receive a control message
		var retryTimeout time.Duration = 0
		var callback chan bool

		// Delay until timeout passes
		timer := time.NewTimer(channel.batchInterval)
	waitLoop:
		for {
			select {
			case event := <-channel.collectChan:
				if event == nil {
					// Channel closed?  Not intercepted by Send()?
					panic("Received nil event")
				}

				buffer = append(buffer, event)
				if len(buffer) >= channel.batchSize {
					break waitLoop
				}

			case ctl := <-channel.controlChan:
				if ctl.stop {
					stopping = true
					if !ctl.flush {
						// No flush? Just exit.
						if ctl.callback != nil {
							channel.signalWhenDone(ctl.callback)
						}
						break mainLoop
					}
				}

				if ctl.flush {
					retryTimeout = ctl.timeout
					callback = ctl.callback
					break waitLoop
				}

			case _ = <-timer.C:
				// Timeout expired
				timer = nil
				break waitLoop
			}
		}

		if timer != nil && !timer.Stop() {
			<-timer.C
		}

		if len(buffer) > 0 {
			// Buffer will be mutated very shortly- capture it before branching
			// of the goroutine to avoid a very real race condition
			go func(buffer TelemetryBufferItems) {
				channel.waitgroup.Add(1)
				defer channel.waitgroup.Done()

				if callback != nil {
					// If we have a callback, wait on the waitgroup now that it's
					// incremented.
					channel.signalWhenDone(callback)
				}

				channel.transmitRetry(buffer, retryTimeout)
			}(buffer)
		} else if callback != nil {
			channel.signalWhenDone(callback)
		}
	}

	close(channel.collectChan)
	close(channel.controlChan)
}

func (channel *InMemoryChannel) transmitRetry(items TelemetryBufferItems, retryTimeout time.Duration) {
	payload := items.serialize()
	retryTimeRemaining := retryTimeout

	for _, wait := range submit_retries {
		result, err := transmit(payload, items, channel.endpointAddress)
		if err == nil && result.IsSuccess() {
			return
		}

		if result.CanRetry() {
			// Filter down to failed items
			payload, items = result.GetRetryItems(payload, items)
			if len(payload) == 0 || len(items) == 0 {
				return
			}
		} else {
			diagnosticsWriter.Write("Telemetry transmission failed; cannot retry\n")
			return
		}

		if retryTimeout > 0 {
			// We're on a time schedule here.  Make sure we don't try longer
			// than we have been allowed.
			if retryTimeRemaining < wait {
				// One more chance left -- we'll wait the max time we can
				// and then retry on the way out.
				time.Sleep(retryTimeRemaining)
				break
			} else {
				// Still have time left to go through the rest of the regular
				// retry schedule
				retryTimeRemaining -= wait
			}
		}

		if result.IsFailure() {
			diagnosticsWriter.Printf("Telemetry transmission failed; retrying in %s\n", wait)
		}

		time.Sleep(wait)
	}

	// One final try
	_, err := transmit(payload, items, channel.endpointAddress)
	if err != nil {
		diagnosticsWriter.Write("Gave up transmitting payload; exhausted retries")
	}
}

func (channel *InMemoryChannel) signalWhenDone(callback chan bool) {
	go func() {
		channel.waitgroup.Wait()
		callback <- true
		close(callback)
	}()
}
