package main 

import (
	"fmt"
	"bufio.Reader"
	"os"
)

var reader Reader

func main() {
	fmt.Println("Pi Alarm Admin Tools. (c) 2018. Kiet T. Tran, Ph.D.")
	reader = bufio.NewReader(os.Stdin)
	if len(os.Args) <= 1 {
		cmdPanel()
	} else {
		cmd := os.Args[1]

		switch cmd {
		case "wifi":
			netName, psk := cfgWifi()

		default:
		}
	}
}

func cmdPanel() {
	fmt.Println("Menu:")
	fmt.Println("1. Configure WiFi Network")
	fmt.Print("\nEnter #:")

	cmd := reader.ReadString('\n')
	switch cmd {
	case "1":
		netName, psk := cfgWifi()
	}
}

func cfgWifi() (string, string) {
	netName := prompt("Enter Network name:")
	psk := prompt("Enter Passkey:")
	fmt.Println(netName, psk)
	return netName, psk
}

func prompt(prompter string) (string) {
	fmt.Print(prompter)
	return reader.ReadString('\n')
}