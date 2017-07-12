package appinsights

import (
	"testing"
	"time"
)

func TestMessageSentToConsumers(t *testing.T) {
	original := "~~~test_message~~~"

	// There may be spurious messages sent by a transmitter's goroutine from another test,
	// so just check that we do get the test message *at some point*.

	listener1chan := make(chan bool)
	listener1 := NewDiagnosticsMessageListener()
	go listener1.ProcessMessages(func(message string) {
		if message == original {
			listener1chan <- true
		}
	})

	listener2chan := make(chan bool)
	listener2 := NewDiagnosticsMessageListener()
	go listener2.ProcessMessages(func(message string) {
		if message == original {
			listener2chan <- true
		}
	})

	diagnosticsWriter.Write(original)

	listener1recvd := false
	listener2recvd := false
	timeout := false
	timer := time.After(time.Second)
	for !(listener1recvd && listener2recvd) && !timeout {
		select {
		case <-listener1chan:
			listener1recvd = true
		case <-listener2chan:
			listener2recvd = true
		case <-timer:
			timeout = true
		}
	}

	if timeout {
		t.Errorf("Message failed to be delivered to both listeners")
	}

	// Clean up
	diagnosticsWriter.listeners = diagnosticsWriter.listeners[:0]
}
