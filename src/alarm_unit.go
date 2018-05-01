package main

import (
 	"encoding/json"
)

type Zone struct {
	ID string			`json:"id,omitempty"`
	Name string			`json:"name,omitempty"`
	
	State string		`json:"state,omitempty"`
	LastReported string	`json:"last-reported,omitempty"`
	LastEvent Event		`json:"last-event,omitempty"`
}

type AlarmUnit struct {
	MACID string 		`json:"macid,omitempty"`
	OwnerName string	`json:"owner-name,omitempty"`
	OwnerEmail string	`json:"owner-email,omitempty"`
	OwnerCell string	`json:"owner-cell,omitempty"`
	NotifyEmail bool	`json:"notify-email,omitempty"`
	NotifyText bool		`json:"notify-text,omitempty"`
	NotifyPhone bool	`json:"notify-phone,omitempty"`
	Zones []Zone		`json:"zones,omitempty"`

	CurrentState string	`json:"current-state,omitempty"`
	DesiredState string `json:"desired-state,omitempty"`
	LastUpdated string 	`json:"last-update,omitempty"`
}

var alarmUnit AlarmUnit

func UnmarshalJsonAlarmUnit(jsonData []byte) (alarmUnit AlarmUnit) {	
    json.Unmarshal(jsonData, &alarmUnit)
    return
}

func MarshalJsonAlarmUnit(alarmUnit AlarmUnit) (jsonData []byte) {
	jsonData, _ = json.Marshal(alarmUnit)
	return
}

func updateAlarmUnitWithEvent(id string, event Event) {
	alarmUnit := getOrMakeAlarmUnitFromDataStore(id)

	switch event.Type {
	case RX_EVENT:
		updateAlarmUnitWithSensor(alarmUnit, UnmarshalJsonSensor([]byte(event.Reason)))
	case REGISTER_EVENT:
		updateAlarmUnitWithRegistration(alarmUnit, UnmarshalJsonAlarmUnit([]byte(event.Reason)))
	case OWNER_EVENT:
	default:
	}
	
	LogMsg("Event processing completed!")
}

func updateAlarmUnitWithSensor(alarmUnit AlarmUnit, sensor Sensor) {

}

func updateAlarmUnitWithRegistration(alarmUnit AlarmUnit, updatedAlarmUnit AlarmUnit) {

}