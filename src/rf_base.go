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
func rfInit(device string, bitRate int) {
	log.Println("device:",device,"-bitRate:",bitRate)

	options := serial.RawOptions
  	options.BitRate = bitRate

  	newPort, err := options.Open(device)

  	if err != nil {
    	log.Panic("severe error", err)
  	}

  	port = newPort
  	log.Println("RF Tx/Rx initialized successfully.")
}

/**
 * receives data from a sensor through rf base station.
 *
 */
func rfReceive() (data string) {
	for {
		buf := make([]byte, 1)
		if c, err := port.Read(buf); err == nil {
			if buf[0] == 'a' {
				buf = make([]byte, 11)
				port.Read(buf)
				data = string(buf)
				log.Println("rfReceive: ", data)
				return
			} else if err != nil {
				log.Println("RfReceive: ERROR!");
   				log.Println(c)
   				log.Panic(err)
   			}
		}
	}

	return
}

func rfReceiver(sensorCh chan Sensor) {
	fmt.Println("\n**** RF Receiver ****\n")
	for {
        data := rfReceive()
        log.Println("managePiAlarm: ", data)
        sensor := makeSensorEvent(data)
        log.Println("**** rfReceiver: ", sensor)
        sensorCh <- sensor
	}
}

