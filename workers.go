package main

import (
	"log"
)

func handleLogs() {
	for {
		msg := <-logChan
		log.Printf("LOG: %s", msg)
	}
}

func handleFederation() {
	for {
		fed := <-fedChan
		fed.Federate()
	}
}
