package main

import (
	"encoding/json"
)

type Network struct {
	WifiSSID string
	WifiPSK string
	WifiKeyManagement string
}

func UnmarshalJsonNetwork(jsonData []byte) (network Network) {	
    json.Unmarshal(jsonData, &network)
    return
}

func MarshalJsonNetwork(alarmUnit Network) (jsonData []byte) {
	jsonData, _ = json.Marshal(network)
	return
}