package main

import (
	"log"
)

func handleLogs() {
	for {
		msg := <-logChan
		log.Printf("LOG: %s", msg)
		// fmt.Printf("log: %s", msg)
	}
}

func handleFederation() {
	for {
		fed := <-fedChan
		// fmt.Println("federating...")
		fed.Federate()
	}
}
