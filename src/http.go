package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "log"
    "time"
    "bytes"
    "strings"
    "fmt"
    "io/ioutil"
    "github.com/mikepb/go-serial"
)

type Event struct {
    ID      string  `json:"id,omitempty"`
    Time 	string  `json:"time,omitempty"`
    Message string  `json:"message,omitempty"`
    Reason  string 	`json:"reason,omitempty"`
}

/**
 * An HTTP Server receives a POST.
 */
func PostEvent(w http.ResponseWriter, r *http.Request) {
	LogMsg("PostEvent: ")

    params := mux.Vars(r)
    var event Event
    _ = json.NewDecoder(r.Body).Decode(&event)

    LogMsg("id: " + params["id"])

    event.Time = time.Now().String()
    event.Reason += " - " + "Received OK!"

    json.NewEncoder(w).Encode(event)
    LogMsg("PostEvent: ends")
}

/**
 * Posting an event to an HTTP Server with a JSON formatted message.
 *
 */
func PostToHttpServer(serverEndpoint string, reason string) {
	LogMsg("PostToHttpServer: " + serverEndpoint)
	LogMsg("PostToHttpServer: " + reason)

	var event Event
		
	event.Time = time.Now().String()
	event.Reason = reason
	event.Message = "An event from PI Alarm"
	event.ID = GetMacAddr()

	jsonBuf, _ := json.Marshal(event)

	serverEndpoint = fmt.Sprintf("%s/event/%s", serverEndpoint, event.ID)

	resp, err := http.Post(serverEndpoint, "application/json", bytes.NewBuffer(jsonBuf))
	if err != nil {
		log.Fatal(err)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		LogMsg("PostToHttpServer: response: " + string(data))
	}

	LogMsg("PostToHttpServer: ends")
}

/**
 */
func ServeHttpProcessEvent(serverEndpoint string) {
	LogMsg("ServeHttpProcessEvent: " + serverEndpoint)
    router := mux.NewRouter()
    router.HandleFunc("/event/{id}", PostEvent).Methods("POST")
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
func ServeRfRxPostHttp(serverEndpoint string) {
	LogMsg("ServeRfRxPostHttp: " + serverEndpoint);

  	options := serial.RawOptions
  	options.BitRate = 9600
  	p, err := options.Open("/dev/ttyAMA0")

  	if err != nil {
    	log.Panic(err)
    	fmt.Println(err)
  	}

  	defer p.Close()
  
	for {
  		buf := make([]byte, 1)
  		if c, err := p.Read(buf); err == nil {
			if buf[0] == 'a' {
				buf = make([]byte, 11)
				p.Read(buf)
				PostToHttpServer(serverEndpoint, string(buf))
			} 
			LogMsg("ServeRfRxPostHttp: '" + string(buf) + "'")
  		} else {
			LogMsg("ServeRfRxPostHttp: ERROR!")
    		log.Println(c)
    		log.Panic(err)
  		}
	}

	LogMsg("ServeRfRxPostHttp: ends");
}
