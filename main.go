package main

import (
  	"github.com/tarm/serial"
  	"log"
	"fmt"
 	"os"
	"net"
)

func main() {
	log.Println("PI Alarm Application ....");
	serverEndpoint := os.Getenv("PI_ALARM_SERVER_ENDPOINT")
	runningAs := os.Getenv("PI_ALARM_RUNNING_MODE")
	fmt.Println("running as: " + runningAs)
	fmt.Println("server endpoint: " + serverEndpoint)

	if "CLIENT" == runningAs {
		processRFReceiver(serverEndpoint)
	} else {
		processAlarmSensors(serverEndpoint)
	}
}

/**
 * This code can run either on a pi or on another host
 * It receives sensor data, process it and store it in a datastore.
 * In addition, it logs all incoming data in the order it receives
 * into a logfile.
 */
func processAlarmSensors(serverEndpoint string) {
	log.Println ("processing alarm sensors begins ... " + serverEndpoint)

	for {
		buf := receiveFromClient(serverEndpoint)
		fmt.Println("'" + buf + "'")
	}
}

/**
 * This code should run on the pi with an attached Wireless Based Station
 * Transmitter/Receiver: 
 * https://ha.privateeyepi.com/store/index.php?route=product/product&product_id=66
 *
 * Main function:
 * Read data from an RF receiver and retransmit exactly to a UDP listener located
 * at PI_ALARM_SERVER_ENDPOINT.
 *
 */
func processRFReceiver(serverEndpoint string) {
	log.Println("processing RF Receiver begins .... " + serverEndpoint);
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600}
	p, err := serial.OpenPort(c)

  	if err != nil {
    	log.Fatal(err)
  	}

  	defer p.Close()
  
	for {
  		buf := make([]byte, 1)
  		if n, err := p.Read(buf); err == nil {
			if buf[0] == 'a' {
				buf = make([]byte, 11)
				p.Read(buf)
				postToServer(serverEndpoint, string(buf))
   				fmt.Println(">>>" + string(buf) + "<<<")
			} else {
    			fmt.Print(buf)
    			fmt.Print(string(buf))
			}
  		} else {
    		log.Println(n)
    		log.Panic(err)
			fmt.Println("PI Alarm Receiver ERROR!....");
  		}
	}

	log.Println("processing RF Receiver ends ....");
}

/**
 * posting buffer to server
 *
 */
func postToServer(serverEndpoint string, buf string) {
	log.Println("posting data to server endpoint ... " + serverEndpoint)
	conn, err := net.Dial("udp", serverEndpoint)	
	if err != nil {
		log.Panic(err)
		return 
	}

	defer conn.Close()
	
	conn.Write([]byte(buf))

	log.Println("posting data ends ...")
}

/**
 * reading from a post from a client
 */
func receiveFromClient(serverEndpoint string) (buf string) {
	conn, err := net.ListenPacket("udp", serverEndpoint)
	
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	buffer := make([] byte, 1024)
	conn.ReadFrom(buffer)

	return string(buffer)
}