// Start a RESTful server to expose total system boot time
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const appVersion string = "0.1"
const httpEndpoint string = "localhost:8080"
const durationPath string = "/duration"
const versionPath string = "/version"
const useJsonOption string = "--json"

//const debugOption string "--debug"
var totalBootTime float64 = -1.0
var useJson bool = false

// getTotalBootTime reads from systemd-analyze system boot time,
// returning a float64 time expressed in seconds
func getTotalBootTime() float64 {
	cmdOutput, err := exec.Command("systemd-analyze", "--system").Output()
	if err != nil {
		log.Fatal(err)
	}
	cmdOutputSplit := strings.Split(string(cmdOutput), "=")
	timeString := strings.TrimSpace(strings.Split(cmdOutputSplit[1], "s\n")[0])
	timeFloat, err := strconv.ParseFloat(timeString, 64)
	if err != nil {
		log.Fatal(err)
	}
	return timeFloat
}

// durationHandler answers requests on durationPath
func durationHandler(writer http.ResponseWriter, request *http.Request) {
	// Check boot time is available
	if totalBootTime <= 0 {
		log.Fatalf("Requested boot time not available, closing web server")
	}
	// Set answer in plain text or json format
	var bootTimeAnswer []byte
	if useJson {
		var bootTime struct {
			Duration float64 `json:"boot_time"`
		}

		bootTime.Duration = totalBootTime
		bootTimeAnswer, _ = json.Marshal(bootTime)
	} else {
		bootTimeAnswer = []byte(fmt.Sprintf("%f", totalBootTime))
	}
	_, err := writer.Write(bootTimeAnswer)
	if err != nil {
		log.Fatal(err)
	}
}

func versionHandler(writer http.ResponseWriter, request *http.Request) {
	var versionAppAnswer []byte
	if useJson {
		var versionApp struct {
			Version string `json:"version"`
		}
		versionApp.Version = appVersion
		versionAppAnswer, _ = json.Marshal(versionApp)
	} else {
		versionAppAnswer = []byte(appVersion)
	}
	_, err := writer.Write(versionAppAnswer)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	totalBootTime = getTotalBootTime()
}

func main() {
	// Do parameter parsing
	if len(os.Args) > 1 {
		if os.Args[1] == useJsonOption && len(os.Args) == 2 {
			useJson = true
		} else {
			fmt.Printf("Usage: %s <%s>", os.Args[0], useJsonOption)
			log.Fatal("Incorrect parameters")
		}
	}

	// Register http handlers
	http.HandleFunc(durationPath, durationHandler)
	http.HandleFunc(versionPath, versionHandler)

	// Start web service
	fmt.Println("Starting server on http://", httpEndpoint)
	err := http.ListenAndServe(httpEndpoint, nil)
	log.Fatal(err)
}
