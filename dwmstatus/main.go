// A simple program for generating the dwm status bar using
// github.com/thomas11/dwmstatus. It outputs the current time and the battery
// level (Linux only).
package main

import (
	"bytes"
	"fmt"
	"github.com/thomas11/dwmstatus"
	"io/ioutil"
	"strconv"
	"time"
)

const (
	sep        = " | "
	timeLayout = "Mon Jan 2 15:04"

	batteryPath = "/sys/class/power_supply/BAT1/"
	batteryNow  = batteryPath + "energy_now"
	batteryFull = batteryPath + "energy_full"
)

// Read an int from a file, assuming that is all the file contains, modulo
// whitespace.
func readIntFromFile(path string) (int, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading from %v: %v\n", path, err)
		return -1, err
	}

	num, err := strconv.Atoi(string(bytes.TrimSpace(content)))
	if err != nil {
		fmt.Printf("Nonsense reading from %v: %v\n", path, err)
		return -1, err
	}
	return num, nil
}

func getBatteryPercentage() int {
	now, err := readIntFromFile(batteryNow)
	if err != nil {
		return -1
	}
	full, err := readIntFromFile(batteryFull)
	if err != nil {
		return -1
	}

	return int(100 * (float32(now) / float32(full)))
}

// This function implements dwmstatus.genTitleFunc.
func genTitle(now time.Time, b *bytes.Buffer) {
	b.WriteString(now.Format(timeLayout))

	b.WriteString(sep)
	b.WriteString(strconv.Itoa(getBatteryPercentage()))
	b.WriteString("%")
}

func main() {
	dwmstatus.Run(3*time.Second, genTitle)
}
