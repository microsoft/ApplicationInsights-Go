# Microsoft Application Insights SDK for Go

[![Build Status](https://travis-ci.org/Microsoft/ApplicationInsights-Go.svg?branch=master)](https://travis-ci.org/Microsoft/ApplicationInsights-Go)

This project provides a Go SDK for Application Insights. [Application Insights](http://azure.microsoft.com/en-us/services/application-insights/) is a service that allows developers to keep their applications available, performant, and successful. This go package will allow you to send telemetry of various kinds (event, metric, trace) to the Application Insights service where they can be visualized in the Azure Portal. 

## Requirements ##
**Install**
```
go get github.com/Microsoft/ApplicationInsights-Go/appinsights
```
**Get an instrumentation key**
>**Note**: an instrumentation key is required before any data can be sent. Please see the "[Getting an Application Insights Instrumentation Key](https://github.com/Microsoft/AppInsights-Home/wiki#getting-an-application-insights-instrumentation-key)" section of the wiki for more information. To try the SDK without an instrumentation key, set the instrumentationKey config value to a non-empty string.

## Usage ##

```go
import "github.com/Microsoft/ApplicationInsights-Go/appinsights"

client := appinsights.NewTelemetryClient("<instrumentation key>")
client.TrackEvent("custom event")
client.TrackMetric("custom metric", 123)
client.TrackTrace("trace message")
```

## telpad test and example app ##

The telpad app can be used to send test telemetry and as an example of using the SDK.

```bash
cd src/github.com/Microsoft/ApplicationInsights-Go
go install

telpad
```