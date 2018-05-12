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
func RfReceive() (data string) {
	buf := make([]byte, 1)
	if c, err := port.Read(buf); err == nil {
		if buf[0] == 'a' {
			buf = make([]byte, 11)
			port.Read(buf)
			log.Println("received: ", (string(buf)))
		} else {
			log.Println("RfReceive: ERROR!");
   			log.Println(c)
   			log.Panic(err)
   		}
	}

	return
}