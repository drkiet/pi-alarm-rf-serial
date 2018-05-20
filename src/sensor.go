package main

import (
	"strings"
	"log"
	"encoding/json"
	"time"
)

const (
	Button = "BUTTON"
	Battery = "BATT"
)

const (
	SENSOR_OPEN = "OPEN"
	SENSOR_CLOSED = "CLOSED"
	SENSOR_NOSTATE = "NOSTATE"
)
type Sensor struct {
	Id string       `json:"id"`
	Type string     `json:"type"`
	ZoneName string `json:"zonename"`
	State string    `json:"state"`
	Subunit string  `json:"subunit"`
	Battery string  `json:"batter"`
	Data string 	`json:"thedata"`
	Updated time.Time `json:"updatedby"`
}

func unmarshalSensor(jsonData []byte) (sensor Sensor) {	
    json.Unmarshal(jsonData, &sensor)
    return
}

func marshalSensor(sensor Sensor) (jsonData []byte) {
	jsonData, _ = json.Marshal(sensor)
	return
}

/**
 * This function processes the buffer receiving from the message from 
 * a wireless sensor (switch, temperature, humidity, camera, motion, flood etc.)
 * and turns it into a Sensor object.
 * 
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
func makeSensorEvent(data string) (sensor Sensor) {
	if len(data) < 11 {
		log.Println("BAD data received", data, "!")
		return 
	}

	sensor.Data = data
	sensor.Id = data[0:2]
	sensor.ZoneName = lookupZoneName(sensor.Id)
	sensor.Updated = time.Now()

	data = data[2:]

	if isButton(data) {
		sensor.Subunit = string(data[6:7])
		sensor.State = string(data[7:9])
		if sensor.State == "ON" {
			sensor.State = SENSOR_CLOSED
		} else if sensor.State == "OF" {
			sensor.State = SENSOR_OPEN
		} else { 
			sensor.State = SENSOR_NOSTATE
		}

		sensor.Type = Button
	} else if strings.HasPrefix(data, Battery) {
		sensor.Type = Battery
	} else {
		log.Println("NOT supported feature: " + data)
	}
	return
}

func isButton(data string) (isButton bool) {
	if strings.HasPrefix(data, Button) {
		isButton = true
	} else {
		isButton = false
	}

	return
}

func isBattery(data string) (isBattery bool) {
	if strings.HasPrefix(data, Battery) {
		isBattery = true
	} else {
		isBattery = false
	}

	return	
}