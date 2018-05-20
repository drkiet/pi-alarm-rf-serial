package main

import (
	"fmt"
	"time"
)

var trackedZones map[string]*Zone

func trackZone(zone *Zone) {
	trackedZones[zone.Id] = zone
}

func isTrackedZone(zone *Zone) (tracked bool) {
	if trackedZones[zone.Id] == nil {
		return false
	} else {
		return true
	}
}

func untrackZone(zone *Zone) {
	trackedZones[zone.Id] = nil
}

func notifyViaEmail (zone *Zone) {
	for name, email := range getToList() {
		subject := fmt.Sprintf("%s.%s: %s", zone.Id, zone.ZoneName, zone.State)
		sendEmail(email, subject, subject)
		fmt.Println("\nSent email to: ", name, "with subject: ", subject)
	}
}

func actNow(zones map[string]*Zone) {
	for _, zone := range zones {
		if zone.State == SENSOR_OPEN {
			if !isTrackedZone(zone) {
				trackZone(zone)
				notifyViaEmail(zone)
			}
		} else if zone.State == SENSOR_CLOSED {
			if isTrackedZone(zone) {
				untrackZone(zone)
				notifyViaEmail(zone)
			}
		} else {
			// fmt.Println(zone.ZoneName, ":", zone.State)
		}
	}
}

func healthMonitor() {
	fmt.Println("\n**** Health Monitor ****")
	lastUpdated := getLastPiAlarmUpdated()
	trackedZones = make(map[string]*Zone)

	// if something change within a sleep of 1 second, act on it.
	for {
		time.Sleep(1000 * time.Millisecond)	
		fmt.Print(".")

		if lastUpdated != getLastPiAlarmUpdated() {
			lastUpdated = getLastPiAlarmUpdated()
			// printZones(getZones())
			actNow(getZones())
		}
		
	}
}