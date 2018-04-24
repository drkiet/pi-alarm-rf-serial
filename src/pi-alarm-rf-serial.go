package main

import (
  	"log"
 	"os"
)

var serverEndpoint, runningAs, repeaterEndpoint, logFileName string
var file os.File

/**
 * Collecting required parameters via OS environment variables.
 * By using environment variables, it is easier to pass into a Docker container
 * and other similar environments.
 */
func init() {
	serverEndpoint = os.Getenv("PI_ALARM_SERVER_ENDPOINT")
	repeaterEndpoint = os.Getenv("PI_ALARM_REPEATER_ENDPOINT")
	runningAs = os.Getenv("PI_ALARM_RUNNING_MODE")
	logFileName = os.Getenv("PI_ALARM_LOG_FILE_NAME")

	file, err := os.OpenFile(logFileName, 
							 os.O_RDWR | 
							 os.O_CREATE | 
							 os.O_APPEND, 0666)
	
	
	LogMsg("*** running as: " + runningAs)
	LogMsg("*** server endpoint: " + serverEndpoint)
	LogMsg("*** repeater endpoint: " + repeaterEndpoint)
	LogMsg("*** Log File Name: " + logFileName)

	if err != nil {
    	log.Fatalf("error opening file: %v", err)
	}
	
	log.SetOutput(file)
}

/**
 * This program can run in multiple modes. Think of it as a small independently 
 * executed application that has a very specific responsibility.
 * 
 * The running mode is specified as a environment variable PI_ALARM_RUNNING_MODE.
 * 
 * 1. Running on a PI as a RF Receiver:
 *    Input: Short (11 bytes or so) message from the RF Tx/Rx board.
 *    Output: an Event message encapsulates the message from the RF Tx/Rx board.
 *	     RF_RECEIVER_TO_HTTP:
 *       Example:
 *           export PI_ALARM_RUNNING_MODE=RF_RECEIVER_TO_HTTP
 *           export PI_ALARM_SERVER_ENDPOINT=192.168.1.63:9090
 *
 *       RF_RECEIVER_TO_UDP:
 *       Example:
 *           export PI_ALARM_RUNNING_MODE=RF_RECEIVER_TO_UDP
 *           export PI_ALARM_SERVER_ENDPOINT=192.168.1.63:9999
 *
 * 2. Running anywhere (PI/Linux/Windows/Docker/OpenShift/etc.) as an event processor
 *    Input: an incoming Event message from UDP or HTTP. The message is processed
 *           here.
 *       EVENT_UDP_SERVER:
 *       Example:
 *           export PI_ALARM_RUNNING_MODE=EVENT_UDP_SERVER
 *           export PI_ALARM_SERVER_ENDPOINT=192.168.1.63:9999
 *
 *       EVENT_HTTP_SERVER:
 *       Example:
 *           export PI_ALARM_RUNNING_MODE=EVENT_HTTP_SERVER
 *           export PI_ALARM_SERVER_ENDPOINT=192.168.1.63:9090
 *
 *    Output: various output conditions based on various incoming events.
 *
 * 3. Running anywhere (PI/Linux/Windows/Docker/OpenShift/etc.) as a UDP Repeater
 *    Input: an incoming Event message from UDP. The message is immediately 
 *           forwarded to a destination without any further processing.
 *    Output: same as inputput to a UDP listener
 *       UDP_REPEATER:
 *       Example:
 *           export PI_ALARM_RUNNING_MODE=UDP_REPEATER
 *           export PI_ALARM_SERVER_ENDPOINT=192.168.1.63:9999
 *           export PI_ALARM_REPEATER_ENDPOINT=172.17.0.2:9999
 *
 */
func main() {

	LogMsg("main: running")

	if "RF_RECEIVER_TO_UDP" == runningAs {
		ServeRfRxPostUdp(serverEndpoint)
	} else if "RF_RECEIVER_TO_HTTP" == runningAs {
		ServeRfRxPostHttp(serverEndpoint)
	} else if "EVENT_UDP_SERVER" == runningAs {
		ServeUdpProcessEvent(serverEndpoint)
	} else if "UDP_REPEATER" == runningAs {
		ServeUdpPostUdp(serverEndpoint, repeaterEndpoint)
	} else if "EVENT_HTTP_SERVER" == runningAs {
		ServeRfRxPostHttp(serverEndpoint)
	} else {
		LogMsg(runningAs + " is an invalid running mode")
	}

	file.Close()
}