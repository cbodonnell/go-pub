package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// --- Configuration --- //

var config Configuration

func getConfig(ENV string) Configuration {
	file, err := os.Open(fmt.Sprintf("config.%s.json", ENV))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config Configuration
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	if config.SSLCert == "" {
		config.Protocol = "http"
	} else {
		config.Protocol = "https"
	}
	return config
}
