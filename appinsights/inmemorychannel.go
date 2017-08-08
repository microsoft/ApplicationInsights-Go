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
	throttle        *throttleManager
}

type inMemoryChannelControl struct {
	// If true, flush the buffer.
	flush bool

	// If true, stop listening on the channel.  (Flush is required if any events are to be sent)
	stop bool

	// If stopping and flushing, this specifies whether to retry submissions on error.
	retry bool

	// If retrying, what is the max time to wait before finishing up?
	timeout time.Duration

	// If specified, a message will be sent on this channel when all pending telemetry items have been submitted
	callback chan bool
}

func NewInMemoryChannel(config *TelemetryConfiguration) *InMemoryChannel {
	channel := &InMemoryChannel{
		endpointAddress: config.EndpointUrl,
		collectChan:     make(chan Telemetry),
		controlChan:     make(chan *inMemoryChannelControl),
		batchSize:       config.MaxBatchSize,
		batchInterval:   config.MaxBatchInterval,
		throttle:        newThrottleManager(),
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

func (channel *InMemoryChannel) IsThrottled() bool {
	return channel.throttle.IsThrottled()
}

func (channel *InMemoryChannel) Close(flush bool, retry bool, timeout time.Duration) chan bool {
	callback := make(chan bool)

	channel.controlChan <- &inMemoryChannelControl{
		stop:     true,
		flush:    flush,
		timeout:  timeout,
		retry:    retry,
		callback: callback,
	}

	return callback
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

			channel.signalWhenDone(ctl.callback)
		}

		if len(buffer) == 0 {
			continue
		}

		// Things that are used by the sender if we receive a control message
		var retryTimeout time.Duration = 0
		var retry bool
		var callback chan bool

		// Delay until timeout passes or buffer fills up
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
					retry = ctl.retry
					if !ctl.flush {
						// No flush? Just exit.
						channel.signalWhenDone(ctl.callback)
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

		// Hold up transmission if we're being throttled
		if !stopping && channel.throttle.IsThrottled() {
			// Channel is currently throttled.  Once the buffer fills, messages will
			// be lost...  If we're exiting, then we'll just try to submit anyway.  That
			// request may be throttled and transmitRetry will perform the backoff correctly.

			diagnosticsWriter.Write("Channel is throttled, events may be dropped.")
			throttleDone := channel.throttle.NotifyWhenReady()
			dropped := 0

		throttledLoop:
			for {
				select {
				case <-throttleDone:
					close(throttleDone)
					break throttledLoop

				case event := <-channel.collectChan:
					// If there's still room in the buffer, then go ahead and add it.
					if len(buffer) < channel.batchSize {
						buffer = append(buffer, event)
					} else {
						if dropped == 0 {
							diagnosticsWriter.Write("Buffer is full, dropping further events.")
						}

						dropped++
					}

				case ctl := <-channel.controlChan:
					if ctl.stop {
						stopping = true
						retry = ctl.retry
						if !ctl.flush {
							channel.signalWhenDone(ctl.callback)
							break mainLoop
						} else {
							// Make an exception when stopping
							break throttledLoop
						}
					}

					// Cannot flush
					// TODO: Figure out what to do about callback?
					if ctl.flush {
						channel.signalWhenDone(ctl.callback)
					}
				}
			}

			diagnosticsWriter.Printf("Channel dropped %d events while throttled", dropped)
		}

		// Send
		if len(buffer) > 0 {
			go func(buffer TelemetryBufferItems, callback chan bool, retry bool, retryTimeout time.Duration) {
				channel.waitgroup.Add(1)
				defer channel.waitgroup.Done()

				if callback != nil {
					// If we have a callback, wait on the waitgroup now that it's
					// incremented.
					channel.signalWhenDone(callback)
				}

				channel.transmitRetry(buffer, retry, retryTimeout)
			}(buffer, callback, retry, retryTimeout)
		} else if callback != nil {
			channel.signalWhenDone(callback)
		}
	}

	close(channel.collectChan)
	close(channel.controlChan)
	channel.throttle.Stop()
}

func (channel *InMemoryChannel) transmitRetry(items TelemetryBufferItems, retry bool, retryTimeout time.Duration) {
	payload := items.serialize()
	retryTimeRemaining := retryTimeout

	for _, wait := range submit_retries {
		result, err := transmit(payload, items, channel.endpointAddress)
		if err == nil && result != nil && result.IsSuccess() {
			return
		}

		if !retry {
			diagnosticsWriter.Write("Refusing to retry telemetry submission (retry==false)")
			return
		}

		// Check for success, determine if we need to retry anything
		if result != nil {
			if result.CanRetry() {
				// Filter down to failed items
				payload, items = result.GetRetryItems(payload, items)
				if len(payload) == 0 || len(items) == 0 {
					return
				}
			} else {
				diagnosticsWriter.Write("Cannot retry telemetry submission")
				return
			}

			// Check for throttling
			if result.IsThrottled() {
				if result.retryAfter != nil {
					diagnosticsWriter.Printf("Channel is throttled until %s", *result.retryAfter)
					channel.throttle.RetryAfter(*result.retryAfter)
				} else {
					// TODO: Pick a time
				}
			}
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

		diagnosticsWriter.Printf("Waiting %s to retry submission", wait)
		time.Sleep(wait)

		// Wait if the channel is throttled and we're not on a schedule
		if channel.IsThrottled() && retryTimeout == 0 {
			diagnosticsWriter.Printf("Channel is throttled; extending wait time.")
			ch := channel.throttle.NotifyWhenReady()
			result := <-ch
			close(ch)

			if !result {
				return
			}
		}
	}

	// One final try
	_, err := transmit(payload, items, channel.endpointAddress)
	if err != nil {
		diagnosticsWriter.Write("Gave up transmitting payload; exhausted retries")
	}
}

func (channel *InMemoryChannel) signalWhenDone(callback chan bool) {
	if callback != nil {
		go func() {
			channel.waitgroup.Wait()
			callback <- true
			close(callback)
		}()
	}
}
