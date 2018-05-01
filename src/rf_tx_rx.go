package main

import (
	"github.com/mikepb/go-serial"
	"log"
	"fmt"
)

var port *serial.Port

/**
 * Create a serial port for use in reading message ...
 * device name: /dev/ttyAMA0
 * speed: 9600
 */
func RfInitialize(device string, bitRate int) {
	options := serial.RawOptions
  	options.BitRate = bitRate

  	newPort, err := options.Open(device)

  	if err != nil {
    	log.Panic(err)
  	}

  	port = newPort
  	LogMsg("RF Tx/Rx initialized successfully.")
}

/**
 * Port must be first initilized. Then, it can receive sensor data/events.
 *
 */
func RfReceive(alarmUnit *AlarmUnit) (sensor Sensor) {
	buf := make([]byte, 1)
	if c, err := port.Read(buf); err == nil {
		if buf[0] == 'a' {
			buf = make([]byte, 11)
			port.Read(buf)
			sensor = transformSensorMessage(alarmUnit, string(buf))
			LogMsg(fmt.Sprintf("RfReceive: ", sensor))
		} else {
			LogMsg("RfReceive: ERROR!");
   			log.Println(c)
   			log.Panic(err)
   		}
	}

	return
}