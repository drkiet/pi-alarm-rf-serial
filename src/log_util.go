package main

import (
	"os"
)


func makeLogFile(fileName string) (file *os.File) {
	file, _ = os.OpenFile(fileName, 
					   os.O_RDWR | 
					   os.O_CREATE | 
					   os.O_APPEND, 0666)
	return
}

