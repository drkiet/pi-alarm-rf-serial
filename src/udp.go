package main
	
import (
	"encoding/json"
	"log"
	"fmt"
	"net"
	"time"
)

/**
 * posting buffer to server
 *
 */
func postToUdpServer(serverEndpoint string, buf string) {
	logMsg("posting data to UDP server endpoint ... " + serverEndpoint)

	conn, err := net.Dial("udp", serverEndpoint)	
	if err != nil {
		log.Panic(err)
		fmt.Println(err)
		return 
	}

	defer conn.Close()
	var event Event
		
	event.Time = time.Now().String()
	event.Reason = buf
	event.Message = "from PI Alarm"
	event.ID = getMacAddr()

	json.NewEncoder(conn).Encode(event)

	logMsg("posting data ends ...")
}