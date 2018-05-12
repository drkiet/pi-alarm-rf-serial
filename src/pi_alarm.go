package main 

import (
	"log"
)

const RFBaseStationSerial = "/dev/ttyAMA0" // best to get it from environment var.
const SerialPortSpeed int = 9600

func managePiAlarm() {
	log.Println("Managing PI Alarm System")

	RfInitialize(RFBaseStationSerial, 9600)

	for {
        data := RfReceive()
        log.Println("managePiAlarm: ", data
	}
}