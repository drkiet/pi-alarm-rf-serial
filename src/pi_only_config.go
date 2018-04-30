package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"io"
)


const ALARM_UNIT_CONFIG = "alarm_unit.properties"
const MAX_ZONES = 20

var alarmUnit AlarmUnit
var network Network


func loadPiAlarmConfigFromFile() {
	file, err := os.Open(configFolder + ALARM_UNIT_CONFIG)
	defer file.Close()
	if err != nil {
		fmt.Println("Need this file: ", configFolder + ALARM_UNIT_CONFIG)
	}

	reader := bufio.NewReader(bufio.NewReader(file))

	alarmUnit.MACID = GetMacAddr()
	alarmUnit.Zones = make([]Zone, 0, MAX_ZONES)
	var zoneIndex int = 0

	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				fmt.Println("ReadLine failed: ", err)
			}
			break
		}
		parseConfig(string(bytes), &zoneIndex)
	}

	fmt.Println("alarmUnit: ", string(MarshalJsonAlarmUnit(alarmUnit)))
	fmt.Println("network: ", string(MarshalJsonNetwork(network)))
}	

func parseConfig(line string, zoneIndex *int) {
	tokens := strings.Split(line, ":")

	switch strings.Trim(tokens[0], " ") {
	case "owner-name":
		alarmUnit.OwnerName = strings.Trim(tokens[1], " ")

	case "owner-email":
		alarmUnit.OwnerEmail = strings.Trim(tokens[1], " ")

	case "owner-cell":
		alarmUnit.OwnerCell = strings.Trim(tokens[1], " ")

	case "notify-email":
		alarmUnit.NotifyEmail = makeYesNoBool(tokens[1])

	case "notify-text":
		alarmUnit.NotifyText = makeYesNoBool(tokens[1])

	case "notify-phone":
		alarmUnit.NotifyPhone = makeYesNoBool(tokens[1])

	case "zone":
		zoneTokens := strings.Split(strings.Trim(tokens[1], " "), "=")
		alarmUnit.Zones = alarmUnit.Zones[:*zoneIndex+1]
		alarmUnit.Zones [*zoneIndex].ID = strings.Trim(zoneTokens[0], " ")
		alarmUnit.Zones [*zoneIndex].Name = strings.Trim(zoneTokens[1], " ")
		*zoneIndex++

	case "wifi-ssid":
		network.WifiSSID = strings.Trim(tokens[1], " ")

	case "wifi-psk":
		network.WifiPSK = strings.Trim(tokens[1], " ")

	case "wifi-key-management":
		network.WifiKeyManagement = strings.Trim(tokens[1], " ")

	default:
		// ignore everything else.
	}	
}

func makeYesNoBool (yesNo string) (trueFalse bool) {
	if strings.ToLower(strings.Trim(yesNo, " ")) == "yes" {
		trueFalse = true
	} else {
		trueFalse = false
	}
	return
}