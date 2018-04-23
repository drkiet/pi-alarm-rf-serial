package main

import (
  	"github.com/mikepb/go-serial"
  	"log"
 	"os"
	"fmt"
	"encoding/json"
	"strings"
)

var serverEndpoint, runningAs, repeaterEndpoint string
var file os.File

func main() {

	log.Println("PI Alarm Application ....");

	if "UDP_CLIENT" == runningAs {
		processUdpClient(serverEndpoint)
	} else if "HTTP_CLIENT" == runningAs {
		processHttpClient(serverEndpoint)
	} else if "UDP_SERVER" == runningAs {
		processUdpServer(serverEndpoint)
	} else if "UDP_REPEATER" == runningAs {
		processUdpRepeater(serverEndpoint, repeaterEndpoint)
	} else if "HTTP_SERVER" == runningAs {
		ProcessHttpServer(serverEndpoint, file)
	} else {
		logMsg("Invalid runningAs: " + runningAs)
	}

	file.Close()
}

func init() {
	file, err := os.OpenFile("./alarm.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	
	serverEndpoint = os.Getenv("PI_ALARM_SERVER_ENDPOINT")
	repeaterEndpoint = os.Getenv("PI_ALARM_REPEATER_ENDPOINT")
	runningAs = os.Getenv("PI_ALARM_RUNNING_MODE")
	
	logMsg("running as: " + runningAs)
	logMsg("server endpoint: " + serverEndpoint)
	logMsg("repeater endpoint: " + repeaterEndpoint)

	if err != nil {
    	log.Fatalf("error opening file: %v", err)
    	fmt.Println(err)
	}
	
	log.SetOutput(file)
}

/**
 * Log into a log file -
 * Then, print it on screen
 */
func logMsg(msg string) {
	log.Println(msg)
	fmt.Println(msg)
}

/**
 * This code can run either on a pi or on another host
 * It receives sensor data, process it and store it in a datastore.
 * In addition, it logs all incoming data in the order it receives
 * into a logfile.
 */
func processUdpServer(serverEndpoint string) {
	log.Println ("processing alarm sensors begins ... " + serverEndpoint)

	for {
		buf := receiveFromUdpClient(serverEndpoint)
		logMsg("received: '" + buf + "'")
	}
}

func ProcessHttpClient(serverEndpoint string) {
	logMsg("processing HTTP Client begins .... " + serverEndpoint);

  	options := serial.RawOptions
  	options.BitRate = 9600
  	p, err := options.Open("/dev/ttyAMA0")

  	if err != nil {
    	log.Panic(err)
    	fmt.Println(err)
  	}

  	defer p.Close()
  
	for {
  		buf := make([]byte, 1)
  		if c, err := p.Read(buf); err == nil {
			if buf[0] == 'a' {
				buf = make([]byte, 11)
				p.Read(buf)
				postToHttpServer(serverEndpoint, string(buf))
			} 
			logMsg("'" + string(buf) + "'")
  		} else {
    		log.Println(c)
    		log.Panic(err)
			log.Println("PI Alarm Receiver ERROR!....");
  		}
	}

	logMsg("processing HTTP Client ends ....");
}
/**
 * This function processes the buffer receiving from the alarm sensor:
 * Starts with:
 * - BUTTON
 * - BTN
 * - TMP
 * - HUM
 * - BATT
 *
 * - SLEEPING
 * - STARTED
 * - AWAKE
 */
func processSensorMessage(buf string) {
	if strings.HasPrefix(buf, "BUTTON") {
		processButton(buf[7:8], buf[8:])
	} else if strings.HasPrefix(buf, "BATT") {
		processBattery()
	} else {
		logMsg("NOT supported feature: " + buf)
	}
}


func processButton(id string, status string) {

}

func processBattery() {

}

/**
 * This function acts as a repeater for the UDP protocol. It receives a message
 * from a server endpoint, the repeate the same message to the repeater endpoint.
 * 
 */
func processUdpRepeater(serverEndpoint string, repeaterEndpoing string) {
	log.Println ("processing UDP Repeater begins ... " + serverEndpoint + " --> " + repeaterEndpoing)

	for {
		buf := receiveFromUdpClient(serverEndpoint)
		logMsg("received: '" + buf + "'")
		var event Event
    	_ = json.Unmarshal([]byte(buf), &event)
		postToUdpServer(repeaterEndpoing, event.Reason)
		logMsg("repeating: '" + event.Reason + "'")
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
func processUdpClient(serverEndpoint string) {
	logMsg("processing UDP Client begins .... " + serverEndpoint);

  	options := serial.RawOptions
  	options.BitRate = 9600
  	p, err := options.Open("/dev/ttyAMA0")

  	if err != nil {
    	log.Panic(err)
    	fmt.Println(err)
  	}

  	defer p.Close()
  
	for {
  		buf := make([]byte, 1)
  		if c, err := p.Read(buf); err == nil {
			if buf[0] == 'a' {
				buf = make([]byte, 11)
				p.Read(buf)
				postToUdpServer(serverEndpoint, string(buf))
			} 
			logMsg("'" + string(buf) + "'")
  		} else {
    		log.Println(c)
    		log.Panic(err)
			log.Println("PI Alarm Receiver ERROR!....");
  		}
	}

	logMsg("processing UDP Client ends ....");
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
func processHttpClient(serverEndpoint string) {
	logMsg("processing HTTP Client begins .... " + serverEndpoint);

  	options := serial.RawOptions
  	options.BitRate = 9600
  	p, err := options.Open("/dev/ttyAMA0")

  	if err != nil {
    	log.Panic(err)
    	fmt.Println(err)
  	}

  	defer p.Close()
  
	for {
  		buf := make([]byte, 1)
  		if c, err := p.Read(buf); err == nil {
			if buf[0] == 'a' {
				buf = make([]byte, 11)
				p.Read(buf)
				postToHttpServer(serverEndpoint, string(buf))
			} 
			logMsg("'" + string(buf) + "'")
  		} else {
    		log.Println(c)
    		log.Panic(err)
			log.Println("PI Alarm Receiver ERROR!....");
  		}
	}

	logMsg("processing HTTP Client ends ....");
}