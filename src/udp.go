package main
	
import (
	"encoding/json"
	"net"
	"time"
	"log"
)

/**
 * posting buffer to server
 *
 */
func PostToUdpServer(serverEndpoint string, reason string) {
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

	json.NewEncoder(conn).Encode(event)
	LogMsg("PostToUdpServer: ends")
}


/**
 * Listening to a UDP connection request & then read the message into a buffer
 */
func ReceiveFromUdpClient(serverEndpoint string) (buf string) {
	LogMsg("ReceiveFromUdpClient: " + serverEndpoint)
	conn, err := net.ListenPacket("udp", serverEndpoint)
	
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	buffer := make([] byte, 1024)
	size, _, _ := conn.ReadFrom(buffer)

	buf = string(buffer[:size])
	LogMsg("ReceiveFromUdpClient: ends")
	return
}

/**
 * This code can run either on a pi or on another host
 * It receives sensor data, process it and store it in a datastore.
 * In addition, it logs all incoming data in the order it receives
 * into a logfile.
 */
func ServeUdpProcessEvent(serverEndpoint string) {
	LogMsg ("ServeUdpProcessEvent: " + serverEndpoint)

	for {
		buf := ReceiveFromUdpClient(serverEndpoint)
		LogMsg("received: '" + buf + "'")
	}
}

/**
 * This function acts as a repeater for the UDP protocol. It receives a message
 * from a server endpoint, the repeate the same message to the repeater endpoint.
 * 
 */
func ServeUdpPostUdp(serverEndpoint string, repeaterEndpoing string) {
	LogMsg ("ServeUdpPostUdp: " + serverEndpoint + " --> " + repeaterEndpoing)

	for {
		buf := ReceiveFromUdpClient(serverEndpoint)
		LogMsg("ServeUdpPostUdp: received: '" + buf + "'")
		var event Event
    	_ = json.Unmarshal([]byte(buf), &event)
		PostToUdpServer(repeaterEndpoing, event.Reason)
		LogMsg("ServeUdpPostUdp: repeating: '" + event.Reason + "'")
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
func ServeRfRxPostUdp(serverEndpoint string) {
	LogMsg("ServeRfRxPostUdp: " + serverEndpoint);

	RfInitialize("/dev/ttyAMA0", 9600)
  
	for {
		sensorEvent := RfReceive()
		PostToUdpServer(serverEndpoint, sensorEvent)
		LogMsg("ServeRfRxPostUdp: '" + sensorEvent + "'")
	}

	LogMsg("ServeRfRxPostUdp: ends");
}

