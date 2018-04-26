package main

import (
 	"encoding/json"
)

type AlarmUnit struct {
	MACID string 		`json:"macid,omitempty"`
	OwnerName string	`json:"owner-name,omitempty"`
	OwnerEmail string	`json:"owner-email,omitempty"`
	OwnerCell string	`json:"owner-cell,omitempty"`
	Zones [] string		`json:"zones,omitempty"`
	CurrentState string	`json:"current-state,omitempty"`
	DesiredState string `json:"desired-state,omitempty"`
	LastUpdate string 	`json:"last-update,omitempty"`
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