package main

import (
	"strings"
	"log"
	"encoding/json"
)

const (
	Button = "BUTTON"
	Battery = "BATT"
)
type Sensor struct {
	SensorId, Type, ZoneName, State, Subunit, Batter string 	
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
	sensor.SensorId = data[0:2]
	sensor.ZoneName = "*** Unknown sensor ***"
	
	data = data[2:]

	if isButton(data) {
		sensor.Subunit = string(data[6:7])
		sensor.State = string(data[7:9])
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