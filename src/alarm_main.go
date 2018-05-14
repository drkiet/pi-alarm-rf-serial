package main

import (
  	"log"
 	"os"
 	"fmt"
)

var runMode, logName string

const (
	ConfigFolder = "./config/"
	LogsFolder = "./logs/"
	LogFile = "./logs/alarm.log"
	EventLogFile = "./logs/events.log"
	PiConfigFile = "./config/pi_alarm.cfg"
)

const (
	RunOnPi = "on_pi"
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

 */
func main() {
	fmt.Println("\n**** Alarm starts ****\n")
	if runMode == RunOnPi {
		managePiAlarm()
	} 


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