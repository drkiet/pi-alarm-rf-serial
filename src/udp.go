package main
	
import (
	"net"
	"fmt"
	"strings"
)

/**
 * posting buffer to server
 *
 */
func postEventToUdpServer(endpoint, id string, event Event) {
	conn, _ := net.Dial("udp", endpoint)	

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
func receiveFromUdpClient(endpoint string) (bufstr string, address string) {
	conn, _ := net.ListenPacket("udp", endpoint)
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
func serveUdpProcessEvent() {
	LogMsg ("ServeUdpProcessEvent: " + serverEndpoint)

    MakeEventStore()
   	getAllAlarmUnits()


	for {
		bufstr, _ := receiveFromUdpClient(serverEndpoint)

		index := strings.Index(bufstr, ";")
		id := bufstr[:index-1]
		jsonEvent := bufstr[index+2:]

		QueueJsonEvent(id, []byte(jsonEvent))
    	event := UnmarshalJsonEvent([]byte(jsonEvent))
    	updateAlarmUnitWithEvent(id, event)
	}
}

/**
 * This function acts as a repeater for the UDP protocol. It receives a message
 * from a server endpoint, the repeate the same message to the repeater endpoint.
 * 
 */
func serveUdpPostUdp() {
	LogMsg ("ServeUdpPostUdp: " + serverEndpoint + " --> " + repeaterEndpoint)

	for {
		bufstr, address := receiveFromUdpClient(serverEndpoint)

		conn, _ := net.Dial("udp", repeaterEndpoint)	
		defer conn.Close()

		bytesWritten, _ := conn.Write([]byte (bufstr))

		LogMsg(fmt.Sprintf("forwarded %s(%d bytes) to %s from %s", 
						   bufstr, bytesWritten, repeaterEndpoint, address))
	}
}

/**
 * This function acts as a forwarder for the UDP protocol to a HTTP protocol 
 * endpoint. It receives a message.
 * 
 */
func serveUdpPostHttp() {
	LogMsg ("ServeUdpPostHttp: " + serverEndpoint + " --> " + repeaterEndpoint)

	for {
		bufstr, address := receiveFromUdpClient(serverEndpoint)

		index := strings.Index(bufstr, ";")
		id := bufstr[:index-1]
		jsonEvent := bufstr[index+2:]

		QueueJsonEvent(id, []byte(jsonEvent))
    	event := UnmarshalJsonEvent([]byte(jsonEvent))
    	response := postEventToHttpServer(repeaterEndpoint, id, event)

		LogMsg(fmt.Sprintf("forwarded %s to %s from %s with response %s", 
						   bufstr, repeaterEndpoint, address, response))
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
func serveRfRxPostUdp() {
	LogMsg("ServeRfRxPostUdp: " + serverEndpoint);

	loadPiAlarmConfigFromFile()
	RfInitialize("/dev/ttyAMA0", 9600)
    alarmUnit := loadPiAlarmConfigFromFile()
    registerThisAlarmUnitWithUdpServer(alarmUnit)

	for {
        postSensorEventToUdpServer(alarmUnit)
	}
}

// Make an RX Event Sensor
// Post the event to the UDP Server
func postSensorEventToUdpServer(alarmUnit AlarmUnit) {
	event := makeSensorEvent(GetMacAddr(), TYPE_RX_EVENT, 
                             RfReceive(&alarmUnit), "from sensor")
	postEventToUdpServer(serverEndpoint, event.ID, event)
}

// Make a registration event for the alarm unit
// Post the event to the UDP Server
func registerThisAlarmUnitWithUdpServer(alarmUnit AlarmUnit) {
    event := makeRegisterEvent(GetMacAddr(), TYPE_REGISTER_EVENT, 
                               alarmUnit, "from host")
    postEventToUdpServer(serverEndpoint, event.ID, event)
}
