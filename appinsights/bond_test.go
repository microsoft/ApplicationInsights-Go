package appinsights

import "testing"

func TestMessageData(t *testing.T) {
	testMessage := "test"

	messageData := &messageData{
		Message: testMessage,
	}

	if messageData.Message != testMessage {
		t.Errorf("Message is %s, want %s", messageData.Message, testMessage)
	}

	messageData.Ver = 2
}
