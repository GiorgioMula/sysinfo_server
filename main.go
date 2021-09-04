// Start a RESTful server to expose total system boot time
package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// getTotalBootTime reads from systemd-analyze system boot time,
// returning a float64 useful
func getTotalBootTime() (float64, error) {
	out, err := exec.Command("systemd-analyze", "--system").Output()
	if err != nil {
		log.Fatal(err)
	}
	out_split := strings.Split(string(out), "=")
	time_string := strings.TrimSpace(strings.Split(out_split[1], "s\n")[0])
	return strconv.ParseFloat(time_string, 64)
}

func main() {
	total_time, _ := getTotalBootTime()
	fmt.Println("System boot total time:", total_time)
}
