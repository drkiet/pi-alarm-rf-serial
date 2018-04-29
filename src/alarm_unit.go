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
	Zones []Zone		`json:"zones,omitempty"`
	CurrentState string	`json:"current-state,omitempty"`
	DesiredState string `json:"desired-state,omitempty"`
	LastUpdated string 	`json:"last-update,omitempty"`
	NotifyEmail bool	`json:"notify-email,omitempty"`
	NotifyText bool		`json:"notify-text,omitempty"`
	NotifyPhone bool	`json:"notify-phone,omitempty"`

}

func UnmarshalJsonAlarmUnit(jsonData []byte) (alarmUnit AlarmUnit) {	
    json.Unmarshal(jsonData, &alarmUnit)
    return
}

func MarshalJsonAlarmUnit(alarmUnit AlarmUnit) (jsonData []byte) {
	jsonData, _ = json.Marshal(alarmUnit)
	return
}