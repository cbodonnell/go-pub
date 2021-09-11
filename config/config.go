package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var (
	defaults = map[string]interface{}{
		"DEBUG":                true,
		"PORT":                 80,
		"LOG_FILE":             "",
		"SERVER_NAME":          "localhost",
		"AUTH":                 "http://localhost:8080/auth",
		"CLIENT":               "http://localhost:3000",
		"ENDPOINTS_USERS":      "users",
		"ENDPOINTS_ACTIVITIES": "activities",
		"ENDPOINTS_OBJECTS":    "objects",
		"ENDPOINTS_INBOX":      "inbox",
		"ENDPOINTS_OUTBOX":     "outbox",
		"ENDPOINTS_FOLLOWING":  "following",
		"ENDPOINTS_FOLLOWERS":  "followers",
		"ENDPOINTS_LIKED":      "liked",
		"SSL_CERT":             "",
		"SSL_KEY":              "",
		"DB_HOST":              "host",
		"DB_PORT":              5432,
		"DB_NAME":              "database",
		"DB_USER":              "user",
		"DB_PASSWORD":          "password",
		"JWT_KEY":              "secret",
		"RSA_PUBLIC_KEY":       "public.pem",
		"RSA_PRIVATE_KEY":      "private.pem",
	}
	configPaths = []string{
		".",
	}
)

// --- Configuration --- //

// Configuration struct
type Configuration struct {
	Debug         bool   `mapstructure:"DEBUG"`
	Port          int    `mapstructure:"PORT"`
	LogFile       string `mapstructure:"LOG_FILE"`
	Protocol      string
	ServerName    string     `mapstructure:"SERVER_NAME"`
	Auth          string     `mapstructure:"AUTH"`
	Client        string     `mapstructure:"CLIENT"`
	Endpoints     Endpoints  `mapstructure:",squash"`
	SSLCert       string     `mapstructure:"SSL_CERT"`
	SSLKey        string     `mapstructure:"SSL_KEY"`
	Db            DataSource `mapstructure:",squash"`
	JWTKey        string     `mapstructure:"JWT_KEY"`
	RSAPublicKey  string     `mapstructure:"RSA_PUBLIC_KEY"`
	RSAPrivateKey string     `mapstructure:"RSA_PRIVATE_KEY"`
}

// DataSource struct
type Endpoints struct {
	Users      string `mapstructure:"ENDPOINT_USERS"`
	Activities string `mapstructure:"ENDPOINT_ACTIVITIES"`
	Objects    string `mapstructure:"ENDPOINT_OBJECTS"`
	Inbox      string `mapstructure:"ENDPOINT_INBOX"`
	Outbox     string `mapstructure:"ENDPOINT_OUTBOX"`
	Following  string `mapstructure:"ENDPOINT_FOLLOWERS"`
	Followers  string `mapstructure:"ENDPOINT_FOLLOWING"`
	Liked      string `mapstructure:"ENDPOINT_LIKED"`
}

// DataSource struct
type DataSource struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	Dbname   string `mapstructure:"DB_NAME"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
}

var (
	C Configuration
)

// TODO: Have defaults for all config variables
func ReadConfig(ENV string) (Configuration, error) {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	viper.SetConfigName(ENV)
	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
	viper.AutomaticEnv()
	var config Configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
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
		return config, err
	}
	config.RSAPrivateKey, err = readKey(config.RSAPrivateKey)
	if err != nil {
		return config, err
	}
	C = config
	return config, nil
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
