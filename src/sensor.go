package main

import (
	"strings"
)

type Sensor struct {
	ID string 		`json:"id,omitempty"`
	Type string 	`json:"type,omitempty"`
	Zone string 	`json:"zone,omitempty"`
	State string 	`json:"state,omitempty"`
	Subunit string 	`json:"subunit,omitempty"`
	Battery string 	`json:"battery,omitempty"`
}
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
func ProcessSensorMessage(buf string) (sensor Sensor) {
	sensor.ID = buf[0:2]
	sensor.Zone = "to-be-looked-up"
	buf = buf[2:]

	if strings.HasPrefix(buf, "BUTTON") {
		sensor.Subunit = string(buf[6:7])
		sensor.State = string(buf[7:9])
		sensor.Type = "BUTTON"
	} else if strings.HasPrefix(buf, "BATT") {
		sensor.Battery = string(buf[4:])
	} else {
		LogMsg("NOT supported feature: " + buf)
	}
	return
}

func ProcessSensorEvent(event Event) {
	LogMsg ("ProcessSensorEvent: " + event.Reason)
	LogMsg ("ProcessSensorEvent: ends")
}