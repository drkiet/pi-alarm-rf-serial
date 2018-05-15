package main

import (
	"fmt"
	"time"
)

func notifyViaEmail (subject string, message string) {
	for name, email := range getToList() {
		sendEmail(email, subject, message + ":" + name)
	}
}

func healthMonitor() {
	fmt.Println("\n**** Health Monitor ****")
	lastUpdated := getLastPiAlarmUpdated()

	// if something change within a sleep of 1 second, act on it.
	for {
		time.Sleep(1000 * time.Millisecond)	
		fmt.Print(".")

		if lastUpdated != getLastPiAlarmUpdated() {
			lastUpdated = getLastPiAlarmUpdated()	
			zonesState := getFormattedZoneStates()
			for _, zoneState := range zonesState {
				subject := zoneState
				notifyViaEmail(subject, subject)
			}
		}
		
	}
}