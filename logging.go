package main

import (
	"fmt"
	"log"
	"os"
)

func setLogFile(file string) *os.File {
	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Logging to %s", file))
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	return logFile
}
