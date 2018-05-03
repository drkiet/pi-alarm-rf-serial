package main

import (
 	"encoding/json"
 	"fmt"
)

type Zone struct {
	ID string			`json:"id,omitempty"`
	Name string			`json:"name,omitempty"`
	
	State string		`json:"state,omitempty"`
	LastReported string	`json:"last-reported,omitempty"`
	LastEvent Event		`json:"last-event,omitempty"`
}

type AlarmUnit struct {
	MACID string 		`json:"macid,omitempty" bson: "macid"`
	OwnerName string	`json:"owner-name,omitempty" bson:"owner-name"`
	OwnerEmail string	`json:"owner-email,omitempty" bson:"owner-email"`
	OwnerCell string	`json:"owner-cell,omitempty" bson:"owner-cell"`
	NotifyEmail bool	`json:"notify-email,omitempty bson:"notify-email"`
	NotifyText bool		`json:"notify-text,omitempty "bson:"notify-text"`
	NotifyPhone bool	`json:"notify-phone,omitempty" bson:"notify-phone"`
	Zones []Zone		`json:"zones,omitempty" bson:"zones"`

	CurrentState string	`json:"current-state,omitempty" bson:"current-state"`
	DesiredState string `json:"desired-state,omitempty bson:"desired-state"`
	LastUpdated string 	`json:"last-update,omitempty" bson:"last-update"`
}

func UnmarshalJsonAlarmUnit(jsonData []byte) (alarmUnit AlarmUnit) {	
    json.Unmarshal(jsonData, &alarmUnit)
    return
}

func MarshalJsonAlarmUnit(alarmUnit AlarmUnit) (jsonData []byte) {
	jsonData, _ = json.Marshal(alarmUnit)
	return
}

func updateAlarmUnitWithEvent(id string, event Event) {
	
	alarmUnit := getAnAlarmUnitFromDataStore(id)

	switch event.Type {
	case TYPE_RX_EVENT:
		if alarmUnit != nil {
			updateAlarmUnitWithSensor(*alarmUnit, event.SensorMsg)
		}

	case TYPE_REGISTER_EVENT:
		if alarmUnit == nil {
			addAnAlarmUnitToDataStore(*alarmUnit)
		}
		updateAlarmUnitWithRegistration(*alarmUnit, event.Alarm)

	case TYPE_OWNER_EVENT:

	default:

	}
	
	LogMsg("Event processing completed!")
}

func updateAlarmUnitWithSensor(alarmUnit AlarmUnit, sensor Sensor) {
	fmt.Println("processing sensor: ", sensor)
}

func updateAlarmUnitWithRegistration(alarmUnit AlarmUnit, updatedAlarmUnit AlarmUnit) {
	fmt.Println("Processing alarm unit: current: ", updatedAlarmUnit, " ---- updated: ", updatedAlarmUnit)
}

// Adding a new Zone to the alarm unit
// This code should be executed only at the host
func addNewZone(alarmUnit *AlarmUnit, zone Zone) {
	fmt.Println(alarmUnit.Zones)
	zoneIndex := len(alarmUnit.Zones)
	alarmUnit.Zones = alarmUnit.Zones[:zoneIndex+1]
	alarmUnit.Zones [zoneIndex].ID = zone.ID
	alarmUnit.Zones [zoneIndex].Name = zone.Name
}