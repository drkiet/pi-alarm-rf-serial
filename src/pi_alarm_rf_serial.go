package main

import (
  	"log"
 	"os"
)

var serverEndpoint, runningAs, repeaterEndpoint, 
	logFileName, configFolder, logsFolder string

var maxBufSize int = 1024


/**
 * Collecting required parameters via OS environment variables.
 * By using environment variables, it is easier to pass into a Docker container
 * and other similar environments.
 */
func init() {
	configFolder = "./config/"
	logsFolder = "./logs/"

	serverEndpoint = os.Getenv("PI_ALARM_SERVER_ENDPOINT")
	repeaterEndpoint = os.Getenv("PI_ALARM_REPEATER_ENDPOINT")
	runningAs = os.Getenv("PI_ALARM_RUNNING_MODE")
	logFileName = os.Getenv("PI_ALARM_LOG_FILE_NAME")

	LogMsg("*** running as: " + runningAs)
	LogMsg("*** server endpoint: " + serverEndpoint)
	LogMsg("*** repeater endpoint: " + repeaterEndpoint)
	LogMsg("*** Log File Name: " + logFileName)

	log.SetOutput(MakeLogFile(logsFolder + logFileName))
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
 *       UDP_UDP_REPEATER:
 *       Example:
 *           export PI_ALARM_RUNNING_MODE=UDP_REPEATER
 *           export PI_ALARM_SERVER_ENDPOINT=192.168.1.63:9999
 *           export PI_ALARM_REPEATER_ENDPOINT=172.17.0.2:9999
 *
 *     	 UDP_HTTP_REPEATER:
 *       Example:
 *           export PI_ALARM_RUNNING_MODE=UDP_HTTP_REPEATER
 *           export PI_ALARM_SERVER_ENDPOINT=192.168.1.63:9999
 *           export PI_ALARM_REPEATER_ENDPOINT=172.17.0.2:9090
 */
func main() {

	LogMsg("main: running")

	if "RF_RECEIVER_TO_UDP" == runningAs {
		ServeRfRxPostUdp()
	} else if "RF_RECEIVER_TO_HTTP" == runningAs {
		ServeRfRxPostHttp()
	} else if "EVENT_UDP_SERVER" == runningAs {
		ServeUdpProcessEvent()
	} else if "UDP_UDP_REPEATER" == runningAs {
		ServeUdpPostUdp()
	} else if "UDP_HTTP_REPEATER" == runningAs {
		ServeUdpPostHttp()
	} else if "EVENT_HTTP_SERVER" == runningAs {
		ServeHttpProcessEvent()
	} else {
		LogMsg(runningAs + " is an invalid running mode")
	}
}