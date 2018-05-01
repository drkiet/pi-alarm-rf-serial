package main

import (
	"strings"
	"encoding/json"
)

type Sensor struct {
	ID string 		`json:"id,omitempty"`
	Type string 	`json:"type,omitempty"`
	Zone string 	`json:"zone,omitempty"`
	State string 	`json:"state,omitempty"`
	Subunit string 	`json:"subunit,omitempty"`
	Battery string 	`json:"battery,omitempty"`
}

func UnmarshalJsonSensor(jsonData []byte) (sensor Sensor) {	
    json.Unmarshal(jsonData, &sensor)
    return
}

func MarshalJsonSensor(sensor Sensor) (jsonData []byte) {
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
func transformSensorMessage(alarmUnit *AlarmUnit, buf string) (sensor Sensor) {
	sensor.ID = buf[0:2]
	sensor.Zone = "*** Unknown sensor ***"
	found := false

	for _, zone := range alarmUnit.Zones {
		if sensor.ID == zone.ID {
			sensor.Zone = zone.Name
			break
		}
	}


	if !found {
		var zone Zone
		zone.ID = sensor.ID
		zone.Name = "*** untracked zone ***"

		addNewZone(alarmUnit, zone)
	}

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
