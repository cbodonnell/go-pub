package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cheebz/go-pub/models"
)

// --- Configuration --- //

var conf models.Configuration

func getConfig(ENV string) models.Configuration {
	// Open config file
	file, err := os.Open(fmt.Sprintf("config.%s.json", ENV))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Decode to Configuration struct
	decoder := json.NewDecoder(file)
	var config models.Configuration
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	// Set Protocol based on SSL config
	if config.SSLCert == "" {
		config.Protocol = "http"
	} else {
		config.Protocol = "https"
	}
	// Read RSA keys
	config.RSAPublicKey, err = readKey(config.RSAPublicKey)
	if err != nil {
		log.Fatal(err)
	}
	config.RSAPrivateKey, err = readKey(config.RSAPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
