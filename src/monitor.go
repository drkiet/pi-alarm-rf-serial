package main

import (
	"fmt"
	"time"
)

func healthMonitor() {
	fmt.Println("\n**** Health Monitor ****")
	lastUpdated := getLastPiAlarmUpdated()

	// if something change within a sleep of 1 second, act on it.
	for {
		time.Sleep(1000 * time.Millisecond)	
		fmt.Print(".")

		if lastUpdated != getLastPiAlarmUpdated() {
			lastUpdated = getLastPiAlarmUpdated()	
			fmt.Println("**** Something changed ...")
		}
		
	}
}