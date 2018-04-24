package main

import (
	"strings"
)

/**
 * This function processes the buffer receiving from the alarm sensor:
 * Starts with:
 * - BUTTON
 * - BTN
 * - TMP
 * - HUM
 * - BATT
 *
 * - SLEEPING
 * - STARTED
 * - AWAKE
 */
func ProcessSensorMessage(buf string) {
	if strings.HasPrefix(buf, "BUTTON") {
		ProcessButton(buf[7:8], buf[8:])
	} else if strings.HasPrefix(buf, "BATT") {
		ProcessBattery()
	} else {
		LogMsg("NOT supported feature: " + buf)
	}
}


func ProcessButton(id string, status string) {

}

func ProcessBattery() {

}

func ProcessSensorEvent(event Event) {
	LogMsg ("ProcessSensorEvent: " + event.Reason)
	LogMsg ("ProcessSensorEvent: ends")
}