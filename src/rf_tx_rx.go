package main

import (
	"github.com/mikepb/go-serial"
	"log"
)

var port *serial.Port

/**
 * Create a serial port for use in reading message ...
 * device name: /dev/ttyAMA0
 * speed: 9600
 */
func RfInitialize(device string, bitRate int) {
	LogMsg("RfInitialize: " + device);
	options := serial.RawOptions
  	options.BitRate = bitRate

  	newPort, err := options.Open(device)

  	if err != nil {
    	log.Panic(err)
  	}

  	port = newPort

  	LogMsg("RfInitialize: ends");
}

/**
 * Port must be first initilized. Then, it can receive sensor data/events.
 *
 */
func RfReceive() (sensorEvent string) {
	LogMsg("RfReceive: " + serverEndpoint);

	buf := make([]byte, 1)
	if c, err := port.Read(buf); err == nil {
		if buf[0] == 'a' {
			buf = make([]byte, 11)
			port.Read(buf)
			sensorEvent = string(buf)
			LogMsg("RfReceive: '" + sensorEvent + "'")
		} else {
			LogMsg("RfReceive: ERROR!");
   			log.Println(c)
   			log.Panic(err)
   		}
	}

	LogMsg("RfReceive: ends");
	return
}