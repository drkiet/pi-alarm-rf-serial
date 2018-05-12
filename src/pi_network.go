package main

import (
	"encoding/json"
	"fmt"
)

type Network struct {
	WifiSSId, WifiPSK, WifiKeyManagement string
}

var network Network

func unmarshalNetwork(jsonData []byte) (network Network) {	
    json.Unmarshal(jsonData, &network)
    return
}

func marshalNetwork(alarmUnit Network) (jsonData []byte) {
	jsonData, _ = json.Marshal(network)
	return
}

func setWifiSSId (wifiSSId string) {
	network.WifiSSId = wifiSSId
}

func setWifiPSK (wifiPSK string) {
	network.WifiPSK = wifiPSK
}

func setWifiKeyManagement (wifiKeyManagement string) {
	network.WifiKeyManagement = wifiKeyManagement
}

func printNetwork() {
	fmt.Println(network)
}