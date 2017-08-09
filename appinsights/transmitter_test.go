package appinsights

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

type testServer struct {
	server *httptest.Server
	notify chan *testRequest

	handler         func(http.ResponseWriter, *http.Request)
	responseData    []byte
	responseCode    int
	responseHeaders map[string]string
}

type testRequest struct {
	request *http.Request
}

func (server *testServer) Close() {
	server.server.Close()
	close(server.notify)
}

func (server *testServer) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if server.handler != nil {
		server.handler(writer, req)
	} else {
		server.defaultHandler(writer, req)
	}

	server.notify <- &testRequest{req}
}

func (server *testServer) defaultHandler(writer http.ResponseWriter, req *http.Request) {
	hdr := writer.Header()

	for k, v := range server.responseHeaders {
		hdr[k] = []string{v}
	}

	writer.WriteHeader(server.responseCode)
	writer.Write(server.responseData)
}

func newTestClientServer() (TelemetryClient, *testServer) {
	server := &testServer{}
	server.server = httptest.NewServer(server)
	server.notify = make(chan *testRequest)
	server.responseCode = 200
	server.responseData = make([]byte, 0)
	server.responseHeaders = make(map[string]string)

	config := NewTelemetryConfiguration("00000000-0000-0000-000000000000")
	config.EndpointUrl = fmt.Sprintf("http://%s/v2/track", server.server.Listener.Addr().String())
	client := NewTelemetryClientFromConfig(config)

	return client, server
}
