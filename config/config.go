package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// --- Configuration --- //

// Configuration struct
type Configuration struct {
	Debug         bool       `json:"debug"`
	Port          int        `json:"port"`
	LogFile       string     `json:"logFile"`
	Protocol      string     `json:"protocol"`
	ServerName    string     `json:"serverName"`
	Auth          string     `json:"auth"`
	Client        string     `json:"client"`
	Endpoints     Endpoints  `json:"endpoints"`
	SSLCert       string     `json:"sslCert"`
	SSLKey        string     `json:"sslKey"`
	Db            DataSource `json:"db"`
	JWTKey        string     `json:"jwtKey"`
	RSAPublicKey  string     `json:"rsaPublicKey"`
	RSAPrivateKey string     `json:"rsaPrivateKey"`
}

// DataSource struct
type Endpoints struct {
	Users      string `json:"users"`
	Activities string `json:"activities"`
	Objects    string `json:"objects"`
	Inbox      string `json:"inbox"`
	Outbox     string `json:"outbox"`
	Following  string `json:"following"`
	Followers  string `json:"followers"`
	Liked      string `json:"liked"`
}

// DataSource struct
type DataSource struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

var (
	C Configuration
)

// TODO: Incorporate dotenv
// TODO: Have defaults for all config variables
func ReadConfig(ENV string) {
	// Open config file
	file, err := os.Open(fmt.Sprintf("config.%s.json", ENV))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Decode to Configuration struct
	decoder := json.NewDecoder(file)
	var config Configuration
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
	C = config
}

func readKey(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return fmt.Sprintf("%s\n", strings.Join(lines, "\n")), nil
}
