package appinsights

type DiagnosticsMessageWriter interface {
	Write(string)
	appendListener(*diagnosticsMessageListener)
}

type diagnosticsMessageWriter struct {
	listeners []chan string
}

type DiagnosticsMessageProcessor func(string)

type DiagnosticsMessageListener interface {
	ProcessMessages(DiagnosticsMessageProcessor)
}

type diagnosticsMessageListener struct {
	channel chan string
}

var writer *diagnosticsMessageWriter = &diagnosticsMessageWriter{
	listeners: make([]chan string, 0),
}

func getDiagnosticsMessageWriter() DiagnosticsMessageWriter {
	return writer
}

func NewDiagnosticsMessageListener() DiagnosticsMessageListener {
	listener := &diagnosticsMessageListener{
		channel: make(chan string),
	}

	writer.appendListener(listener)

	return listener
}

func (writer *diagnosticsMessageWriter) appendListener(listener *diagnosticsMessageListener) {
	writer.listeners = append(writer.listeners, listener.channel)
}

func (writer *diagnosticsMessageWriter) Write(message string) {
	for _, c := range writer.listeners {
		c <- message
	}
}

func (listener *diagnosticsMessageListener) ProcessMessages(process DiagnosticsMessageProcessor) {
	for {
		message := <-listener.channel
		process(message)
	}
}
