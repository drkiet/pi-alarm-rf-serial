package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"io"
)




func loadPiAlarmCfg() {
	file, err := os.Open(PiConfigFile)
	defer file.Close()
	if err != nil {
		fmt.Println("Need this file: ", PiConfigFile)
	}

	reader := bufio.NewReader(bufio.NewReader(file))

	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				fmt.Println("ReadLine failed: ", err)
			}
			break
		}
		parseConfig(string(bytes))
	}

	printPiAlarm()
	printNetwork()

	return
}	

func parseConfig(line string) {
	tokens := strings.Split(line, ":")

	switch strings.Trim(tokens[0], " ") {
	case "owner":
		setOwner(strings.Trim(tokens[1], " "))

	case "email":
		setEmail(strings.Trim(tokens[1], " "))

	case "cell":
		setCell(strings.Trim(tokens[1], " "))

	case "notify-via":
		setNotifyVia(strings.Trim(tokens[1], " "))

	case "zone":
		addZone(strings.Trim(tokens[1], " "))

	case "wifi-ssid":
		setWifiSSId(strings.Trim(tokens[1], " "))

	case "wifi-psk":
		setWifiPSK(strings.Trim(tokens[1], " "))

	case "wifi-key-management":
		setWifiKeyManagement(strings.Trim(tokens[1], " "))

	default:
		// ignore everything else.
	}	
}

