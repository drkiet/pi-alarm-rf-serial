package main

import (
	"bufio"
	"os"
	"fmt"
)

const ALARM_UNIT_CONFIG = "alarm_unit.properties"

func loadPiAlarmConfigFromFile() {
	file, err := os.Open(configFolder + ALARM_UNIT_CONFIG)
	defer file.Close()
	if err != nil {
		fmt.Println("Need this file: ", configFolder + ALARM_UNIT_CONFIG)
	}

	reader := bufio.NewReader(bufio.NewReader(file))

	bytes, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("ReadLine failed: ", err)
	}

	fmt.Println("line: ", string(bytes))
}	