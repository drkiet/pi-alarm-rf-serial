package main
	
import (
	"encoding/json"
	"net"
	"time"
	"log"
	"fmt"
)

/**
 * posting buffer to server
 *
 */
func PostToUdpServer(reason string) {
	LogMsg("PostToUdpServer: " + serverEndpoint)

	conn, err := net.Dial("udp", serverEndpoint)	
	if err != nil {
		log.Fatal(err)
		return 
	}

	defer conn.Close()
	var event Event
		
	event.Time = time.Now().String()
	event.Reason = reason
	event.Message = "from PI Alarm"
	event.ID = GetMacAddr()

	buf, _ := json.Marshal(event)
	LogMsg("PostToUdpServer: posting: " + string(buf))

	json.NewEncoder(conn).Encode(event)
	LogMsg("PostToUdpServer: ends")
}


/**
 * Listening to a UDP connection request & then read the message into a buffer
 */
func ReceiveFromUdpClient() (buf string, address string) {
	LogMsg("ReceiveFromUdpClient: " + serverEndpoint)
	conn, err := net.ListenPacket("udp", serverEndpoint)
	
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	buffer := make([] byte, 1024)
	size, addr, _ := conn.ReadFrom(buffer)

	buf = string(buffer[:size])
	address = addr.String()
	LogMsg("ReceiveFromUdpClient: ends")
	return
}

/**
 * This code can run either on a pi or on another host
 * It receives sensor data, process it and store it in a datastore.
 * In addition, it logs all incoming data in the order it receives
 * into a logfile.
 */
func ServeUdpProcessEvent() {
	LogMsg ("ServeUdpProcessEvent: " + serverEndpoint)

	for {
		buf, address := ReceiveFromUdpClient()
		logtext := fmt.Sprintf("received %s from %s", buf, address)
		LogMsg("ServeUdpProcessEvent: " + logtext)

		var event Event
    	_ = json.Unmarshal([]byte(buf), &event)
		ProcessSensorEvent(event)
	}
}

/**
 * This function acts as a repeater for the UDP protocol. It receives a message
 * from a server endpoint, the repeate the same message to the repeater endpoint.
 * 
 */
func ServeUdpPostUdp() {
	LogMsg ("ServeUdpPostUdp: " + serverEndpoint + " --> " + repeaterEndpoint)

	for {
		buf, address := ReceiveFromUdpClient()
		var event Event
    	_ = json.Unmarshal([]byte(buf), &event)
		PostToUdpServer(event.Reason)

		logtext := fmt.Sprintf("receiving %s from %s and repeating %s to %s", 
							   buf, address, event.Reason, repeaterEndpoint)
		LogMsg("ServeUdpPostUdp: " + logtext)
	}
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
func ServeRfRxPostUdp() {
	LogMsg("ServeRfRxPostUdp: " + serverEndpoint);

	RfInitialize("/dev/ttyAMA0", 9600)
  
	for {
		sensorEvent := RfReceive()
		PostToUdpServer(sensorEvent)
		LogMsg("ServeRfRxPostUdp: '" + sensorEvent + "'")
	}
}

