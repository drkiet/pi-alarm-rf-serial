package main

import (
    "encoding/json"
    "os"
    "fmt"
    "github.com/golang-collections/go-datastructures/queue"
)

type Event struct {
    ID      string  `json:"id,omitempty"`
    Type    string  `json:"type,omitempty"`
    Reason  string 	`json:"reason,omitempty"`
    Time 	string  `json:"time,omitempty"`
    Message string  `json:"message,omitempty"`
}

var eventFile *os.File
var q *queue.Queue

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

func processEvent(id string, event Event) {
	if "RX_EVENT" == event.Type {
		sensor := ProcessSensorMessage(event.Reason)
		msg, _ := json.Marshal(sensor)
		LogMsg("Sesnor: " + string(msg))
	} else {
		LogMsg("Unprocessed Type: " + event.Type)
	}

	LogMsg("Event processing completed!")
}
