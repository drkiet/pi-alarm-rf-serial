package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "log"
    "time"
    "bytes"
    "fmt"
    "io/ioutil"
)

/**
 * An event arrives
 * The event is queued
 * The response is immediate returned to publisher/poster
 * 
 *
 */
func processPostEvent(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    id := params["id"]
    jsonEvent, _ := ioutil.ReadAll(r.Body)

    QueueJsonEvent(id, jsonEvent)
    event := UnmarshalJsonEvent(jsonEvent)
    updateAlarmUnitWithEvent(id, event)

    event.Type = "HTTP_RESPONSE"
    event.Time = time.Now().String()
    event.Message = "Received OK!"

    jsonEvent = MarshalJsonEvent(event)
   	fmt.Fprintf(w, "%s", jsonEvent)

   	QueueJsonEvent(id, jsonEvent)
}

/**
 * Posting an event to an HTTP Server with a JSON formatted message.
 *
 */
func postEventToHttpServer(endpoint, id string, event Event) (responseEvent Event) {
	jsonEvent := MarshalJsonEvent(event)

	httpEndpoint := fmt.Sprintf("%s/event/%s", endpoint, id)
    
    LogMsg("jsonEvent: " + string(jsonEvent))
    LogMsg("httpEndpoint: " + httpEndpoint)

	resp, err := http.Post(httpEndpoint, "application/json", 
				 	     bytes.NewBuffer(jsonEvent))
	LogMsg("Posted: " + string(jsonEvent))
    
    if err != nil {
        fmt.Println("HTTP post error: ", err)
        return
    }
    
    jsonEvent, _ = ioutil.ReadAll(resp.Body)
    responseEvent = UnmarshalJsonEvent(jsonEvent)
    LogMsg("Response: " + string(jsonEvent))
        
	return
}

/**
 */
func serveHttpProcessEvent() {
	LogMsg("ServeHttpProcessEvent: serving " + serverEndpoint)
    
    router := mux.NewRouter()
    router.HandleFunc("/event/{id}", processPostEvent).Methods("POST")

    MakeEventStore()
    getAllAlarmUnits()

    log.Fatal(http.ListenAndServe(serverEndpoint, router))
}



/**
 * This code should run on the pi with an attached Wireless Based Station
 * Transmitter/Receiver: 
 * https://ha.privateeyepi.com/store/index.php?route=product/product&product_id=66
 *
 * Main function:
 * - load the alarm configuration into the AlarmUnit object. Each host(pi) 
 *   is represented by a single AlarmUnit.
 * - register the AlarmUnit with the Server
 * - receives sensor data, process it and post it to the Server
 *
 */
func serveRfRxPostHttp() {
	LogMsg("Serving RF Rx & Http posting")

	RfInitialize("/dev/ttyAMA0", 9600)
    loadPiAlarmConfigFromFile()
    registerThisAlarmUnitWithHttpServer()

	for {
        postSensorEventToHttpServer()
	}
}

// Make a RX Event Sensor
// Post the event to the HTTP Server
func postSensorEventToHttpServer() {
    event := makeSensorEvent(GetMacAddr(), RX_EVENT, 
                             RfReceive(), "from sensor")
    postEventToHttpServer(serverEndpoint, event.ID, event)
}

// Make a registration event for the alarm unit
// Post the event to the HTTP Server
func registerThisAlarmUnitWithHttpServer() {
    event := makeRegisterEvent(GetMacAddr(), REGISTER_EVENT, 
                               alarmUnit, "from host")
    postEventToHttpServer(serverEndpoint, event.ID, event)
}

