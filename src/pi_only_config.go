package main

import (
	"bufio"
	"os"
	"fmt"
	"io"
)

const ALARM_UNIT_CONFIG = "alarm_unit.properties"

func loadPiAlarmConfigFromFile() {
	file, err := os.Open(configFolder + ALARM_UNIT_CONFIG)
	defer file.Close()
	if err != nil {
		fmt.Println("Need this file: ", configFolder + ALARM_UNIT_CONFIG)
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

		fmt.Println("line: ", string(bytes))
	}
}	