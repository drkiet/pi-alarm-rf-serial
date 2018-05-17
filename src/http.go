package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "log"
    "fmt"
    "strings"
)

var httpEndpoint string

func setHttpEndpoint(serverEndpoint string) {
	httpEndpoint = strings.Trim(serverEndpoint, " ")
}

func manageAlarm(w http.ResponseWriter, r *http.Request) {
   	fmt.Fprintf(w, "Hello, world!")
}

/**
 */
func serveHttp() {
	log.Println("serveHttp: serving '" + httpEndpoint + "'")
    
    router := mux.NewRouter()
    router.HandleFunc("/manage", manageAlarm).Methods("GET")

    log.Fatal(http.ListenAndServe("192.168.1.63:9090", router))
}

