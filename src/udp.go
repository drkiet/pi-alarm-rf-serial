package main
	
import (
	"net"
	"fmt"
	"log"
	"bufio"
	"os"
	"strings"
)

var udpReceiverEndpoint string 
var maxBufSize int = 1024

func udpInit(udpEndpoint string) {
	udpReceiverEndpoint = udpEndpoint
}

func udpReceiver(sensorCh chan Sensor) {
	fmt.Println("\n**** UDP Receiver ****\n")
	for {
        data := udpReceive()
        sensor := makeSensorEvent(data)
        
        if sensor.Id != "" {
        	log.Println("processing ", sensor)
        	sensorCh <- sensor
        } else {
        	log.Println("bad data received.")
        }
	}
}

func udpReceive() (data string) {
	conn, _ := net.ListenPacket("udp", udpReceiverEndpoint)
	defer conn.Close()

	buf := make([] byte, maxBufSize)
	size, _, _ := conn.ReadFrom(buf)
	buf = buf[:size]
	data = string(buf)

	return
}

func udpSend(udpEndpoint string, data string) {
	conn, _ := net.Dial("udp", udpEndpoint)	
	defer conn.Close()
	bytesWritten, _ := conn.Write([]byte (data))
	fmt.Println("Write", bytesWritten, "bytes as ", data)
}

func onKeyboard(udpEndpoint string) {
	fmt.Println("\n**** On Keyboard ***\n")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("enter: ")
		text, _ := reader.ReadString('\n')
		udpSend(udpEndpoint, strings.ToUpper(text[:len(text)-1]))
	}
}