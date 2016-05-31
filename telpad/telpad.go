package main

import "bufio"
import "fmt"
import "os"
import "strconv"
import "strings"

import "github.com/Microsoft/ApplicationInsights-Go/appinsights"

const version string = "1.0"

func processDiagnosticMessage(message string) {
	fmt.Println(message)
	fmt.Printf("> ")
}

func main() {
	fmt.Println("telpad " + version + " - a tool for sending test telemetry.")

	iKey := "2fd42776-5762-41fa-b141-d27dfdb31b3f"

	if len(os.Args) == 2 {
		iKey = os.Args[1]
		if len(iKey) < 36 {
			fmt.Printf("Error: Invalid iKey '%s' specified.", iKey)
			fmt.Println()
			return
		}
	}

	telemetryClient := appinsights.NewTelemetryClient(iKey)

	fmt.Println()
	fmt.Printf("Sending telemetry with iKey '%s'.", iKey)

	fmt.Println()

	diagnosticsMessageListener := appinsights.NewDiagnosticsMessageListener()
	go diagnosticsMessageListener.ProcessMessages(processDiagnosticMessage)

	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, _, err := in.ReadLine()
		if err != nil {
			panic(err)
		}

		input := string(line)

		if input == "exit" {
			return
		}

		if input == "" {
			continue
		}

		parts := strings.Split(input, " ")
		numParts := len(parts)
		if numParts >= 2 {
			itemType := parts[0]
			switch itemType {
			case "event":
				telemetryClient.TrackEvent(strings.TrimSpace(input[len(itemType):]))
			case "metric":
				if numParts < 3 {
					fmt.Println("Metric value is required.")
					break
				}
				value, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 32)
				if err != nil {
					fmt.Println("Invalid metric value.")
					break
				}
				telemetryClient.TrackMetric(strings.TrimSpace(parts[1]), float32(value))
			case "request":

			default:
				telemetryClient.TrackTrace(input)
			}
		} else {
			telemetryClient.TrackTrace(input)
		}
	}
}
