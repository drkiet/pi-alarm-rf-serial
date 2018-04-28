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
func postEvent(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    id := params["id"]
    jsonEvent, _ := ioutil.ReadAll(r.Body)

    QueueJsonEvent(id, jsonEvent)
    event := UnmarshalJsonEvent(jsonEvent)
    processEvent(id, event)

    event.Type = "HTTP_RESPONSE"
    event.Time = time.Now().String()
    event.Reason += " - " + "Received OK!"

    jsonEvent = MarshalJsonEvent(event)
   	fmt.Fprintf(w, "%s", jsonEvent)

   	QueueJsonEvent(id, jsonEvent)
}

func postRegistration(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

    id := params["id"]
    jsonRegistration, _ := ioutil.ReadAll(r.Body)
    QueueJsonEvent(id, jsonRegistration)

    var event Event
    event.ID = id
    event.Type = "REGISTRATION_EVENT"
    event.Reason = "REGISTER"
    event.Time = time.Now().String()
    event.Message = string(jsonRegistration)

    processEvent(id, event)
}

/**
 * Posting an event to an HTTP Server with a JSON formatted message.
 *
 */
func PostToHttpServer(endpoint, id string, event Event) (responseEvent Event) {
	jsonEvent := MarshalJsonEvent(event)

	httpEndpoint := fmt.Sprintf("%s/event/%s", endpoint, id)
    
    LogMsg("jsonEvent: " + string(jsonEvent))
    LogMsg("httpEndpoint: " + httpEndpoint)

	resp, _ := http.Post(httpEndpoint, "application/json", 
				 	     bytes.NewBuffer(jsonEvent))
	LogMsg("Posted: " + string(jsonEvent))
	
	jsonEvent, _ = ioutil.ReadAll(resp.Body)
	responseEvent = UnmarshalJsonEvent(jsonEvent)
	LogMsg("Response: " + string(jsonEvent))
	return
}

/**
 */
func ServeHttpProcessEvent() {
	LogMsg("ServeHttpProcessEvent: serving " + serverEndpoint)
    
    router := mux.NewRouter()
    router.HandleFunc("/event/{id}", postEvent).Methods("POST")
    router.HandleFunc("/registration/{id}", postRegistration).Methods("POST")

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
	LogMsg("Serving RF Rx & Http posting")

	RfInitialize("/dev/ttyAMA0", 9600)

	for {
		var event Event
		event.ID = GetMacAddr()
		event.Type = "RX_EVENT"
		event.Reason = RfReceive()
		event.Time = time.Now().String()
		event.Message = "from sensor"
		
		PostToHttpServer(serverEndpoint, event.ID, event)
	}
}