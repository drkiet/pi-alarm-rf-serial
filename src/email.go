package main

import (
	"strings"
	"fmt"
	"io"
	"bufio"
	"os"
)

var routineEmail string
var routinePsw string
var routineTos [] string
var toCounter int = 0

func emailInit() {
	routineTos = make([]string, 0, 10)

	file, err := os.Open(".private")
	defer file.Close()
	if err != nil {
		fmt.Println("Need this file for email/text: .private")
	}

	reader := bufio.NewReader(bufio.NewReader(file))

	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				fmt.Println("ReadLine failed: ", err)
			}
			break
		}
		parseEmailSettings(string(bytes))
	}
}


func parseEmailSettings(line string) {
	tokens := strings.Split(line, ":")

	switch strings.Trim(tokens[0], " ") {
	case "email":
		routineEmail = (strings.Trim(tokens[1], " "))

	case "psw":
		routinePsw = (strings.Trim(tokens[1], " "))

	case "to":
		routineTos = routineTos[:toCounter+1]
		routineTos[toCounter] = (strings.Trim(tokens[1], " "))
		toCounter++

	default:
	}
}
