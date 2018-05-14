package main

import (
	"fmt"
	"time"
)

func healthMonitor() {
	fmt.Println("\n**** Health Monitor ****")
	for {
		time.Sleep(1000 * time.Millisecond)	
		fmt.Print(".")
	}
}