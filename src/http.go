package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "os"
    "log"
    "time"
)

type Event struct {
    ID      string  `json:"id,omitempty"`
    Time 	string  `json:"time,omitempty"`
    Message string  `json:"message,omitempty"`
    Reason  string 	`json:"reason,omitempty"`
}

func PostEvent(w http.ResponseWriter, r *http.Request) {
	logMsg("Posting event ...")

    params := mux.Vars(r)
    var event Event
    _ = json.NewDecoder(r.Body).Decode(&event)

    logMsg(params["id"])
    logMsg("id: " + event.ID)
    logMsg("Time: " + event.Time)
    logMsg("Message: " + event.Message)
    logMsg("Reason: " + event.Reason)
    logMsg("Posting event ends ...")

    event.Time = time.Now().String()
    event.Reason += " - " + "Processed OK!"

    json.NewEncoder(w).Encode(event)
}

func ProcessHttpServer(serverEndpoint string, file os.File) {
	logMsg("Listening on HTTP ... " + serverEndpoint)
    router := mux.NewRouter()
    router.HandleFunc("/event/{id}", PostEvent).Methods("POST")
    log.Fatal(http.ListenAndServe(serverEndpoint, router))
}
