package main

import (
    "encoding/json"
    "os"
    "time"
)

type Event struct {
    SourceId, Type, Data string
    Time                 time.Time
}

var eventFile *os.File

/**
 * unmarshal an event from a json object to a struct
 */
func unmarshalEvent(jsonData []byte) (event Event) {	
    json.Unmarshal(jsonData, &event)
    return
}

/**
 * marshal an event from a struct into a json object
 */
func marshalEvent(event Event) (jsonData []byte) {
	jsonData, _ = json.Marshal(event)
	return
}

/**
 * make event
 */
func makeEvent(sourceId string, eventType string, data string) (event *Event) {
    event = new(Event)
    event.SourceId = sourceId
    event.Type = eventType
    event.Data = data
    event.Time = time.Now()
    return
}

/**
 * create/open an events log to record an event as it arrives
 */
func makeEventLog() {
	eventFile = makeLogFile(EventLogFile)
}

/**
 * append an event json object to event log
 */
func recordEvent(eventType string, data string) {
    event := makeEvent(getMacAddr(), eventType, data)
    eventFile.Write(marshalEvent(*event))
    eventFile.WriteString("\n")
}