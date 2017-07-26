package appinsights

import (
	"time"
)

type throttleManager struct {
	msgs chan *throttleMessage
}

type throttleMessage struct {
	query     bool
	wait      bool
	throttle  bool
	stop      bool
	timestamp time.Time
	result    chan bool
}

func newThrottleManager() *throttleManager {
	result := &throttleManager{
		msgs: make(chan *throttleMessage),
	}

	go result.run()
	return result
}

func (throttle *throttleManager) RetryAfter(t time.Time) {
	throttle.msgs <- &throttleMessage{
		throttle:  true,
		timestamp: t,
	}
}

func (throttle *throttleManager) IsThrottled() bool {
	ch := make(chan bool)
	throttle.msgs <- &throttleMessage{
		query:  true,
		result: ch,
	}

	result := <-ch
	close(ch)
	return result
}

func (throttle *throttleManager) NotifyWhenReady() chan bool {
	result := make(chan bool, 1)
	throttle.msgs <- &throttleMessage{
		wait:   true,
		result: result,
	}

	return result
}

func (throttle *throttleManager) Stop() {
	result := make(chan bool)
	throttle.msgs <- &throttleMessage{
		stop:   true,
		result: result,
	}

	<-result
	close(result)
}

func (throttle *throttleManager) run() {
mainLoop:
	for {
		// --- Not throttled ---
		var throttledUntil time.Time

	notThrottledLoop:
		for {
			msg := <-throttle.msgs
			if msg.query {
				msg.result <- false
			} else if msg.wait {
				msg.result <- true
			} else if msg.stop {
				break mainLoop
			} else if msg.throttle {
				throttledUntil = msg.timestamp
				break notThrottledLoop
			}
		}

		duration := throttledUntil.Sub(time.Now())
		if duration < 0 {
			continue
		}

		var notify []chan bool

		// --- Throttled and waiting ---
		t := time.NewTimer(duration)

	throttleLoop:
		for {
			select {
			case <-t.C:
				for _, n := range notify {
					n <- true
				}

				break throttleLoop
			case msg := <-throttle.msgs:
				if msg.query {
					msg.result <- true
				} else if msg.wait {
					notify = append(notify, msg.result)
				} else if msg.stop {
					for _, n := range notify {
						n <- false
					}

					msg.result <- true

					break mainLoop
				} else if msg.throttle {
					if msg.timestamp.After(throttledUntil) {
						throttledUntil = msg.timestamp

						if !t.Stop() {
							<-t.C
						}

						t.Reset(throttledUntil.Sub(time.Now()))
					}
				}
			}
		}
	}

	close(throttle.msgs)
}
