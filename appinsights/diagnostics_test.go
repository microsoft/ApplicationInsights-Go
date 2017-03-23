package appinsights

import "testing"
import "sync"

func TestMessageSentToConsumers(t *testing.T) {
	original := "test"

	var wg sync.WaitGroup
	wg.Add(2)

	listener1 := NewDiagnosticsMessageListener()
	go listener1.ProcessMessages(func(message string) {
		if message != original {
			t.Errorf("listener1 returned difference messages, want '%s' got '%s'.", original, message)
		}
		wg.Done()
	})

	listener2 := NewDiagnosticsMessageListener()
	go listener2.ProcessMessages(func(message string) {
		if message != original {
			t.Errorf("listener2 returned difference messages, want '%s' got '%s'.", original, message)
		}
		wg.Done()
	})

	diagnosticsWriter.Write(original)

	wg.Wait()
}
