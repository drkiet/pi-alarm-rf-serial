package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "log"
    "strings"
    "flag"
    "fmt"
    "sort"
    "encoding/json"
)

// error response contains everything we need to use http.Error
type handlerError struct {
    Error   error
    Message string
    Code    int
}

var httpEndpoint string

func setHttpEndpoint(serverEndpoint string) {
	httpEndpoint = strings.Trim(serverEndpoint, " ")
}

// a custom type that we can use for handling errors and formatting responses
type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)

// attach the standard ServeHTTP method to our handler so the http library can call it
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // here we could do some prep work before calling the handler if we wanted to

    // call the actual handler
    response, err := fn(w, r)

    // check for errors
    if err != nil {
        log.Printf("ERROR: %v\n", err.Error)
        http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
        return
    }
    if response == nil {
        log.Printf("ERROR: response from method is nil\n")
        http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
        return
    }

    // turn the response into JSON
    bytes, e := json.Marshal(response)
    if e != nil {
        http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
        return
    }

    // send the response and log
    w.Header().Set("Content-Type", "application/json")
    w.Write(bytes)
    log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}

func listZones(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
    zones := make([]Zone, 0)

    var ids [] string
    for _, zone := range getZones() {
        ids = append(ids, zone.Id)
    }
    sort.Strings(ids)

    for _, id := range ids {
        var theZone Zone
        zone := getZones()[id]
        theZone.Id = zone.Id
        theZone.ZoneName = zone.ZoneName
        theZone.State    = zone.State
        zones = append(zones, theZone)
    }
    return zones, nil
}

func getZone(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
    id := mux.Vars(r)["id"]
    for _, zone := range getZones() {
        if id == zone.Id {
            var theZone Zone
            theZone.Id = zone.Id
            theZone.ZoneName = zone.ZoneName
            theZone.State    = zone.State
            return theZone, nil
        }
    }
    return "NOTFOUND", nil
}

/**
 */
func httpServer(serverEndpoint string) {
	log.Println("serveHttp: serving '" + serverEndpoint + "'")
    // command line flags
    // port := flag.Int("port", 9999, "port to serve on")
    dir := flag.String("directory", "web/", "directory of web files")
    flag.Parse()

    // handle all requests by serving a file of the same name
    fs := http.Dir(*dir)
    fileHandler := http.FileServer(fs)

    router := mux.NewRouter()
    router.Handle("/", http.RedirectHandler("/static/", 302))
    router.Handle("/zones", handler(listZones)).Methods("GET")
    // router.Handle("/zones", handler(addZone)).Methods("POST")
    router.Handle("/zones/{id}", handler(getZone)).Methods("GET")
    // router.Handle("/zones/{id}", handler(updateZone)).Methods("POST")
    // router.Handle("/zones/{id}", handler(removeZone)).Methods("DELETE")
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileHandler))
    http.Handle("/", router)

    log.Fatal(http.ListenAndServe(serverEndpoint, router))
}

