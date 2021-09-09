package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func SetLogFile(file string) *os.File {
	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Logging to %s", file))
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	return logFile
}

func LogCaller(err error) {
	_, path, line, _ := runtime.Caller(2)
	file := filepath.Base(path)
	log.Printf("%s:%d: %v", file, line, err)
}
