package main

import (
  	"log"
 	"os"
 	"fmt"
)

var runMode, logName, cfgServerEndpoint, cfgUDPEndpoint string
var piAlarmVersion string = "0.1"

const (
	ConfigFolder = "./config/"
	LogsFolder = "./logs/"
	LogFile = "./logs/alarm.log"
	EventLogFile = "./logs/events.log"
	PiConfigFile = ".private"
)

const (
	RunsOnPi = "on_pi"
)


/**
 * Collecting required parameters via OS environment variables.
 * By using environment variables, it is easier to pass into a Docker container
 * and other similar environments.
 */
func init() {
	makeDir(ConfigFolder)
	makeDir(LogsFolder)
	log.SetOutput(makeLogFile(LogFile))

	runMode = os.Getenv("RUN_MODE")
	cfgServerEndpoint = os.Getenv("SERVER_ENDPOINT")
	cfgUDPEndpoint = os.Getenv("UDP_ENDPOINT")
	
	fmt.Println("run-mode: ", runMode)
	fmt.Println("server-endpoint: ", cfgServerEndpoint)
	fmt.Println("udp-endpoint: ", cfgUDPEndpoint)
}

func runsOnPi() (runsOnPi bool) {
	return runMode == RunsOnPi
}

func getServerEndpoint() (serverEndpoint string) {
	return cfgServerEndpoint
}

func getVersion() (version string) {
	return piAlarmVersion
}

func getUdpEndpoint() (udpEndpoint string) {
	return cfgUDPEndpoint
}

/**
 * Make the directory if not exist
 */
func makeDir(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, os.ModePerm)
	}
}

/**
 * This program can run in multiple modes. Think of it as a small independently 
 * executed application that has a very specific responsibility.
 * 
 * When this program runs 'on_pi', it behaves this way:
 * - loads cfg from .private file
 * - listens to RF base station receiver for sensors messages
 * - monitors the changes in the alarm application
 * 		+ sends out notification to interested parties (i.e., emails/text)
 *		+ sends out notification to a central webservice
 * - listens to http session for the following commands:
 *		+ modify configuration
 *		+ pull the status of the alarm system
 */
func main() {
	fmt.Println("\n**** Alarm starts ****\n")
	
	if runMode != "on_pi" && runMode != "on_udp" {
		
		if runMode == "on_kb" {
			onKeyboard(cfgUDPEndpoint)
		}

		log.Panic("Must specify correct running mode: on_pi or on_udp")
	}

	managePiAlarm()


	// else if "RF_RECEIVER_TO_HTTP" == runningAs {
	// 	serveRfRxPostHttp()
	// } else if "EVENT_UDP_SERVER" == runningAs {
	// 	serveUdpProcessEvent()
	// } else if "UDP_UDP_REPEATER" == runningAs {
	// 	serveUdpPostUdp()
	// } else if "UDP_HTTP_REPEATER" == runningAs {
	// 	serveUdpPostHttp()
	// } else if "EVENT_HTTP_SERVER" == runningAs {
	// 	serveHttpProcessEvent()
	// } else {
	// 	fmt.Println(runningAs + " is an invalid running mode")
	// }
}