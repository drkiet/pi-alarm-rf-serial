package main

import (
  	"github.com/mikepb/go-serial"
  	"log"
	"fmt"
 	"os"
	"net"
)

func main() {
	log.Println("PI Alarm Receiver running....");
	serverEndpoint := os.Getenv("PI_ALARM_SERVER_ENDPOINT")
	fmt.Println("serverEndpoint: " + string(os.Getenv("PI_ALARM_SERVER_ENDPOINT")))
	processRFReceiver(serverEndpoint)
}

func processRFReceiver(serverEndpoint string) {
	log.Println("processing RF Receiver begins .... " + serverEndpoint);

  	options := serial.RawOptions
  	options.BitRate = 9600
  	p, err := options.Open("/dev/ttyAMA0")

  	if err != nil {
    	log.Panic(err)
  	}

  	defer p.Close()
  
	for {
  		buf := make([]byte, 1)
  		if c, err := p.Read(buf); err == nil {
			if buf[0] == 'a' {
				buf = make([]byte, 11)
				p.Read(buf)
				postToServer(serverEndpoint, string(buf))
   				fmt.Println(">>>" + string(buf) + "<<<")
			} else {
    			fmt.Print(buf)
    			fmt.Print(string(buf))
			}
  		} else {
    		log.Println(c)
    		log.Panic(err)
			fmt.Println("PI Alarm Receiver ERROR!....");
  		}
	}

	log.Println("processing RF Receiver ends ....");
}

/**
 * posting buffer to server
 *
 */
func postToServer(serverEndpoint string, buf string) {
	log.Println("posting data to server endpoint ... " + serverEndpoint)
	conn, err := net.Dial("udp", serverEndpoint)	
	if err != nil {
		log.Panic(err)
		return 
	}

	defer conn.Close()
	
	conn.Write([]byte(buf))

	log.Println("posting data ends ...")
}
