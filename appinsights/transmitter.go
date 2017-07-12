package appinsights

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type transmissionResult struct {
	statusCode int
	retryAfter *time.Time
	response   *backendResponse
}

// Structures returned by data collector
type backendResponse struct {
	itemsReceived int
	itemsAccepted int
	errors        []*itemTransmissionResult
}

type itemTransmissionResult struct {
	index      int
	statusCode int
	message    string
}

const (
	successResponse                         = 200
	partialSuccessResponse                  = 206
	requestTimeoutResponse                  = 408
	tooManyRequestsResponse                 = 429
	tooManyRequestsOverExtendedTimeResponse = 439
	errorResponse                           = 500
	serviceUnavailableResponse              = 503
)

func transmit(payload []byte, items TelemetryBufferItems, endpoint string) (*transmissionResult, error) {
	if endpoint == "" {
		// Special case for tests: don't actually send telemetry to empty endpoint address
		diagnosticsWriter.Write("Refusing to transmit telemetry to empty endpoint\n")
		return &transmissionResult{statusCode: successResponse}, nil
	}

	diagnosticsWriter.Printf("\n----------- Transmitting %d items ---------\n\n", len(items))
	startTime := time.Now()

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-json-stream")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		diagnosticsWriter.Printf("Failed to transmit telemetry: %s\n", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime)

	result := &transmissionResult{statusCode: resp.StatusCode}

	// Grab Retry-After header
	if retryAfterValue, ok := resp.Header[http.CanonicalHeaderKey("Retry-After")]; ok && len(retryAfterValue) == 1 {
		if retryAfterTime, err := time.Parse("", retryAfterValue[0]); err != nil {
			result.retryAfter = &retryAfterTime
		}
	}

	// Parse body, if possible
	response := &backendResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		result.response = response
	}

	// Write diagnostics
	if diagnosticsWriter.hasListeners() {
		diagnosticsWriter.Printf("Telemetry transmitted in %s\n", duration)
		diagnosticsWriter.Printf("Response: %d\n", result.statusCode)
		if result.response != nil {
			diagnosticsWriter.Printf("Items accepted/received: %d/%d\n", result.response.itemsAccepted, result.response.itemsReceived)
			if len(result.response.errors) > 0 {
				diagnosticsWriter.Printf("Errors:\n")
				for _, err := range result.response.errors {
					diagnosticsWriter.Printf("#%d - %d %s\n", err.index, err.statusCode, err.message)
					diagnosticsWriter.Printf("Telemetry item:\n\t%s\n", err.index, string(items[err.index:err.index+1].serialize()))
				}
			}
		}
	}

	return result, nil
}

func (result *transmissionResult) IsSuccess() bool {
	return result.statusCode == successResponse ||
		// Partial response but all items accepted
		(result.statusCode == partialSuccessResponse &&
			result.response != nil &&
			result.response.itemsReceived == result.response.itemsAccepted)
}

func (result *transmissionResult) IsFailure() bool {
	return result.statusCode != successResponse && result.statusCode != partialSuccessResponse
}

func (result *transmissionResult) CanRetry() bool {
	return result.statusCode == partialSuccessResponse ||
		(result.retryAfter != nil &&
			(result.statusCode == requestTimeoutResponse ||
				result.statusCode == serviceUnavailableResponse ||
				result.statusCode == errorResponse ||
				result.statusCode == tooManyRequestsResponse ||
				result.statusCode == tooManyRequestsOverExtendedTimeResponse))
}

func (result *transmissionResult) IsPartialSuccess() bool {
	return result.statusCode == partialSuccessResponse &&
		result.response != nil &&
		result.response.itemsReceived != result.response.itemsAccepted
}

func (result *itemTransmissionResult) CanRetry() bool {
	return result.statusCode == requestTimeoutResponse ||
		result.statusCode == serviceUnavailableResponse ||
		result.statusCode == errorResponse ||
		result.statusCode == tooManyRequestsResponse ||
		result.statusCode == tooManyRequestsOverExtendedTimeResponse
}

func (result *transmissionResult) GetRetryItems(payload []byte, items TelemetryBufferItems) ([]byte, TelemetryBufferItems) {
	if result.statusCode == partialSuccessResponse && result.response != nil {
		// Make sure errors are ordered by index
		sort.Slice(result.response.errors, func(i, j int) bool {
			return result.response.errors[i].index < result.response.errors[j].index
		})

		var resultPayload bytes.Buffer
		resultItems := make(TelemetryBufferItems, 0)
		ptr := 0
		idx := 0

		// Find each retryable error
		for _, responseResult := range result.response.errors {
			if responseResult.CanRetry() {
				// Advance ptr to start of desired line
				for ; idx < responseResult.index && ptr < len(payload); ptr++ {
					if payload[idx] == '\n' {
						idx++
					}
				}

				startIdx := idx

				// Read to end of line
				for ; idx == responseResult.index && ptr < len(payload); ptr++ {
					if payload[idx] == '\n' {
						idx++
					}
				}

				// Copy item into output buffer
				resultPayload.Write(payload[startIdx:idx])
				resultItems = append(resultItems, items[responseResult.index])
			}
		}

		return resultPayload.Bytes(), resultItems
	} else if result.CanRetry() {
		return payload, items
	} else {
		return payload[:0], items[:0]
	}
}
