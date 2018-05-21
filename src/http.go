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
    "io/ioutil"
    "errors"
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

func listZonesHandler(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
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
        theZone.Data     = zone.Data
        theZone.Updated  = zone.Updated
        zones = append(zones, theZone)
    }
    return zones, nil
}


func getZoneHandler(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
    id := mux.Vars(r)["id"]
    for _, zone := range getZones() {
        if id == zone.Id {
            var theZone Zone
            theZone.Id = zone.Id
            theZone.ZoneName = zone.ZoneName
            theZone.State    = zone.State
            theZone.Data     = zone.Data
            theZone.Updated  = zone.Updated
            return theZone, nil
        }
    }
    return "NOTFOUND", nil
}

func removeZoneHandler(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
    id := mux.Vars(r)["id"]
    zone := getZoneById(id)

    if zone == nil {
        return Zone{}, &handlerError{errors.New("Zone not exists!"), "Zone not exists", http.StatusBadRequest}
    }

    removeZone(id)

    return "Successfully removed " + id + "!", nil
}

func updateZoneHandler(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
    id := mux.Vars(r)["id"]
    zone := getZoneById(id)

    if zone == nil {
        return Zone{}, &handlerError{errors.New("Zone not exists"), "Zone not exists", http.StatusBadRequest}
    }

    payload, _ := parseZoneRequest(r)
    updateZoneName(payload.Id, payload.ZoneName)
    zone = getZoneById(id)

    return zone, nil
}

func addZoneHandler(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
    fmt.Println("adding zone ...")
    id := mux.Vars(r)["id"]
    zone := getZoneById(id)

    fmt.Println(zone)
    if zone != nil {
        return Zone{}, &handlerError{errors.New("Zone exists"), "Zone exists", http.StatusBadRequest}
    }

    payload, _ := parseZoneRequest(r)
    addZone(payload.Id + "=" + payload.ZoneName)
    return payload, nil
}

func parseZoneRequest(r *http.Request) (Zone, *handlerError) {
    // the zone payload is in the request body
    data, e := ioutil.ReadAll(r.Body)
    if e != nil {
        return Zone{}, &handlerError{e, "Could not read request", http.StatusBadRequest}
    }

    // turn the request body (JSON) into a zone object
    var payload Zone
    e = json.Unmarshal(data, &payload)
    if e != nil {
        return Zone{}, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
    }

    return payload, nil
}

func parseCmdRequest(r *http.Request) (Cmd, *handlerError) {
    // the cmd payload is in the request body
    data, e := ioutil.ReadAll(r.Body)
    if e != nil {
        return Cmd{}, &handlerError{e, "Could not read request", http.StatusBadRequest}
    }

    // turn the request body (JSON) into a zone object
    var payload Cmd
    e = json.Unmarshal(data, &payload)
    if e != nil {
        return Cmd{}, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
    }

    return payload, nil
}

type Cmd struct {
    Passcode string `json:"passcode"`
    Exec string     `json:"exec"`
    Result string   `json:"result"`
}

func cmdHandler(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
    payload, _ := parseCmdRequest(r)

    if isAuthorized(payload.Passcode) {
        switch payload.Exec {
        case "ARM":
            setWantedState(ARMED)
            payload.Result = "SUCCESS"

        case "DISARM":
            setWantedState(DISARMED)
            payload.Result = "SUCCESS"

        case "PERIMETER":
            setWantedState(PERIMETERED)
            payload.Result = "SUCCESS"

        case "AUTHORIZE":
            if isAuthorized(payload.Passcode) {
                payload.Result = "AUTHORIZED"
            } else {
                payload.Result = "UNAUTH"
            }
        default:
            payload.Result = "BADCMD"
        }

    } else {
        payload.Result = "UNAUTH"
    }

    return payload, nil

}    

func getSystemHandler(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
    sysInfos := make([]PiAlarm, 0)
    sysInfo := getSystemInfo()
    sysInfos = append(sysInfos, *sysInfo)
    return sysInfos, nil
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
    router.Handle("/zones", handler(listZonesHandler)).Methods("GET")
    router.Handle("/zones/{id}", handler(addZoneHandler)).Methods("POST")
    router.Handle("/zones/{id}", handler(getZoneHandler)).Methods("GET")
    router.Handle("/zones/{id}", handler(updateZoneHandler)).Methods("PUT")
    router.Handle("/zones/{id}", handler(removeZoneHandler)).Methods("DELETE")

    router.Handle("/cmd", handler(cmdHandler)).Methods("POST")

    router.Handle("/sysinfo", handler(getSystemHandler)).Methods("GET")

    router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileHandler))
    http.Handle("/", router)

    log.Fatal(http.ListenAndServe(serverEndpoint, router))
}

