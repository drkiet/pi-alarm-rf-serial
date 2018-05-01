package main

import (
    "github.com/golang-collections/go-datastructures/queue"
    "encoding/json"
    "os"
    "fmt"
    "time"
)

const (
	RX_EVENT = "RX_EVENT"
	REGISTER_EVENT = "REGISTER_EVENT"
    OWNER_EVENT = "OWNER_EVENT"
)

type Event struct {
    ID          string    `json:"id,omitempty"`
    Type        string    `json:"type,omitempty"`
    Reason      string 	  `json:"reason,omitempty"`
    Alarm       AlarmUnit `json:"alarm,omitempty"`
    SensorMsg   Sensor    `json:"sensor-msg,omitempty"`
    Time 	    string    `json:"time,omitempty"`
    Message     string    `json:"message,omitempty"`
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
	eventFile = MakeLogFile(logsFolder + "events.log")
}



func makeSensorEvent(id string, eventType string, sensor Sensor, msg string) (event Event) {
    event.ID = id
    event.Type = eventType
    event.SensorMsg = sensor
    event.Time = time.Now().String()
    event.Message = msg
    return
}

func makeRegisterEvent(id string, eventType string, alarmUnit AlarmUnit, msg string) (event Event) {
    event.ID = id
    event.Type = eventType
    event.Alarm = alarmUnit
    event.Time = time.Now().String()
    event.Message = msg
    return
}