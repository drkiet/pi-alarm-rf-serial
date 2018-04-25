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
func PostEvent(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]
    
    jsonEvent, _ := ioutil.ReadAll(r.Body)
    QueueJsonEvent(id, jsonEvent)

    event := UnmarshalJsonEvent(jsonEvent)
    event.Type = "HTTP_RESPONSE"
    event.Time = time.Now().String()
    event.Reason += " - " + "Received OK!"

    jsonEvent = MarshalJsonEvent(event)
   	fmt.Fprintf(w, "%s", jsonEvent)
   	QueueJsonEvent(id, jsonEvent)
}

/**
 * Posting an event to an HTTP Server with a JSON formatted message.
 *
 */
func PostToHttpServer(event Event) (response Event) {
	jsonEvent := MarshalJsonEvent(event)

	serverEndpoint = fmt.Sprintf("%s/event/%s", serverEndpoint, event.ID)
	resp, _ := http.Post(serverEndpoint, "application/json", 
						   bytes.NewBuffer(jsonEvent))
	LogMsg("Posted: " + string(jsonEvent))
	
	jsonEvent, _ = ioutil.ReadAll(resp.Body)
	response = UnmarshalJsonEvent(jsonEvent)
	LogMsg("Response: " + string(jsonEvent))
	return
}

/**
 */
func ServeHttpProcessEvent() {
	LogMsg("ServeHttpProcessEvent: serving " + serverEndpoint)
    router := mux.NewRouter()
    router.HandleFunc("/event/{id}", PostEvent).Methods("POST")
    MakeEventStore()
    log.Fatal(http.ListenAndServe(serverEndpoint, router))
}



/**
 * This code should run on the pi with an attached Wireless Based Station
 * Transmitter/Receiver: 
 * https://ha.privateeyepi.com/store/index.php?route=product/product&product_id=66
 *
 * Main function:
 * Read data from an RF receiver and retransmit exactly to a UDP listener located
 * at PI_ALARM_SERVER_ENDPOINT.
 *
 */
func ServeRfRxPostHttp() {
	for {
		var event Event
		event.ID = GetMacAddr()
		event.Type = "RX_EVENT"
		event.Reason = RfReceive()
		event.Time = time.Now().String()
		event.Message = "from sensor"
		
		PostToHttpServer(event)
	}
}