package main
	
import (
	"net"
	"time"
	"fmt"
	"strings"
)

/**
 * posting buffer to server
 *
 */
func PostToUdpServer(id string, event Event) {
	conn, _ := net.Dial("udp", serverEndpoint)	

	defer conn.Close()

	jsonEvent := MarshalJsonEvent(event)
	bufstr := id + "; " + string(jsonEvent)
	bytesWritten, _ := conn.Write([]byte (bufstr))

	LogMsg(fmt.Sprintf("posted %s (%d bytes) via udp: ", 
					   string(jsonEvent), bytesWritten))
}


/**
 * Listening to a UDP connection request & then read the message into a buffer
 */
func ReceiveFromUdpClient() (bufstr string, address string) {
	conn, _ := net.ListenPacket("udp", serverEndpoint)
	defer conn.Close()

	buf := make([] byte, maxBufSize)
	size, addr, _ := conn.ReadFrom(buf)
	buf = buf[:size]
	bufstr = string(buf)
	address = addr.String()

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

    MakeEventStore()

	for {
		bufstr, _ := ReceiveFromUdpClient()

		index := strings.Index(bufstr, ";")
		id := bufstr[:index-1]
		jsonEvent := bufstr[index+2:]

		QueueJsonEvent(id, []byte(jsonEvent))
    	event := UnmarshalJsonEvent([]byte(jsonEvent))
    	processEvent(id, event)
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

		conn, _ := net.Dial("udp", serverEndpoint)	
		defer conn.Close()

		bytesWritten, _ := conn.Write([]byte (buf))

		LogMsg(fmt.Sprintf("forwarded %s(%d bytes) to %s from %s", 
						   string(buf), bytesWritten, repeaterEndpoint, address))
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
		var event Event
		event.ID = GetMacAddr()
		event.Type = "RX_EVENT"
		event.Reason = RfReceive()
		event.Time = time.Now().String()
		event.Message = "from sensor"

		PostToUdpServer(event.ID, event)
	}
}

