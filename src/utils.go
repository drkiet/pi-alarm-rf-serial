package main

import (
	"net"
	"bytes"
	"log"
	"os"
)

/**
 * Getting a MAC address from hardware/virtual.
 */
func GetMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}

/**
 * Log into a log file -
 * Then, print it on screen
 */
func LogMsg(msg string) {
	log.Println(msg)
}


func MakeLogFile(fileName string) (file *os.File) {
	file, _ = os.OpenFile(fileName, 
					   os.O_RDWR | 
					   os.O_CREATE | 
					   os.O_APPEND, 0666)
	return
}

