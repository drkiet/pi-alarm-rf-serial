package main

import (
 	"encoding/json"
 	"fmt"
)

type Zone struct {
	ID string			`json:"id,omitempty" bson:"id"`
	Name string			`json:"name,omitempty" bson:"name"`
	
	State string		`json:"state,omitempty" bson:"state"`
	LastReported string	`json:"last-reported,omitempty" bson:"last-reported"`
	LastEvent Event		`json:"last-event,omitempty" bson:"last-event"`
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
	// dumpAlarmUnit(*alarmUnit)

	switch event.Type {
	case TYPE_RX_EVENT:
		if alarmUnit != nil {
			for _, zone := range alarmUnit.Zones {
				if zone.ID == event.SensorMsg.ID {
					fmt.Println("old state: ", zone.State)
					break
				}
			}
			updateAlarmUnitWithSensor(alarmUnit, event.SensorMsg)
			for _, zone := range alarmUnit.Zones {
				if zone.ID == event.SensorMsg.ID {
					fmt.Println("new state: ", zone.State)
					break
				}
			}
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

func updateAlarmUnitWithSensor(alarmUnit *AlarmUnit, sensor Sensor) {
	var foundZone bool = false
	for i, zone := range alarmUnit.Zones {
		if zone.ID == sensor.ID {
			fmt.Println("*** Found sensor: ", )
			alarmUnit.Zones[i].State = sensor.State
			foundZone = true
			break
		}
	}

	if !foundZone {
		fmt.Println("Error in locating zone for sensor ID: ", sensor.ID)
	}
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

func dumpAlarmUnit(au AlarmUnit) {
	
	
	fmt.Println("         Name: ", au.OwnerName, "\n",
				"        Email: ", au.OwnerEmail, "\n",
				"         Cell: ", au.OwnerCell, "\n",
				" DesiredState: ", au.DesiredState, "\n",
				"Current State: ", au.CurrentState, "\n",
				" Last Updated: ", au.LastUpdated, "\n",
				" Notify Email: ", au.NotifyEmail, "\n",
				"  Notify Text: ", au.NotifyText, "\n",
				" Notify Phone: ", au.NotifyPhone, "\n\n");
	for _, zone := range au.Zones {
		fmt.Println("           ID: ", zone.ID, "\n", 
				    "         Name: ", zone.Name, "\n", 
				    "        State: ", zone.State, "\n",
				    "Last Reported: ", zone.LastReported, "\n\n")
		fmt.Println("ID: ", zone.LastEvent.ID, "\n",
					"Type: ", zone.LastEvent.Type, "\n",
					"Time: ", zone.LastEvent.Time, "\n",
					"Message: ", zone.LastEvent.Message, "\n\n")
	}
}

