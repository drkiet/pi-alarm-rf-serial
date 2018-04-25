package main

import (
    "encoding/json"
    "os"
    "fmt"
)

type Event struct {
    ID      string  `json:"id,omitempty"`
    Type    string  `json:"type,omitempty"`
    Reason  string 	`json:"reason,omitempty"`
    Time 	string  `json:"time,omitempty"`
    Message string  `json:"message,omitempty"`
}

var eventFile *os.File

func UnmarshalJsonEvent(jsonData []byte) (event Event) {	
    json.Unmarshal(jsonData, &event)
    return
}

func MarshalJsonEvent(event Event) (jsonData []byte) {
	jsonData, _ = json.Marshal(event)
	return
}

func QueueJsonEvent(id string, jsonEvent []byte) {
	line2Write := fmt.Sprintf("id: %s; event:%s\n", id, jsonEvent)
	eventFile.Write([]byte(line2Write))
}

func MakeEventStore() {
	eventFile = MakeLogFile("./event.log")
}