package main 

import (

)

// Search for an alarm unit identified by an ID.
// If found, return the alarmunit.
// Else, create the alarm unit and store it inthe datastore
func getAnAlarmUnitFromDataStore(macid string) (alarmUnit *AlarmUnit) {
	alarmUnit = getAnAlarmUnit(macid)
	return
}

func addAnAlarmUnitToDataStore(alarmUnit AlarmUnit) (success bool) {
	return
}

func deleteAnAlarmUnitFromDataStore(id string) (success bool) {
	return
}

func updateAnAlarmUnitInDataStore(id string, alarmUnit AlarmUnit) (success bool) {
	return
}

func getAllAlarmUnitsFromDataStore() (alarmUnits []AlarmUnit) {
	return getAllAlarmUnits()
}