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
		"DEBUG":               true,
		"PORT":                80,
		"LOG_FILE":            "",
		"SERVER_NAME":         "localhost",
		"AUTH":                "http://localhost:8080/auth",
		"CLIENT":              "http://localhost:3000",
		"ENDPOINT_USERS":      "users",
		"ENDPOINT_ACTIVITIES": "activities",
		"ENDPOINT_OBJECTS":    "objects",
		"ENDPOINT_INBOX":      "inbox",
		"ENDPOINT_OUTBOX":     "outbox",
		"ENDPOINT_FOLLOWING":  "following",
		"ENDPOINT_FOLLOWERS":  "followers",
		"ENDPOINT_LIKED":      "liked",
		"SSL_CERT":            "",
		"SSL_KEY":             "",
		"DB_HOST":             "host",
		"DB_PORT":             5432,
		"DB_NAME":             "database",
		"DB_USER":             "user",
		"DB_PASSWORD":         "password",
		"JWT_KEY":             "secret",
		"RSA_PUBLIC_KEY":      "public.pem",
		"RSA_PRIVATE_KEY":     "private.pem",
		"REDIS_HOST":          "localhost",
		"REDIS_PORT":          6379,
		"REDIS_PASSWORD":      "",
		"REDIS_DB":            0,
		"REDIS_EXP_SECONDS":   3600,
		"ALLOWED_ORIGINS":     "",
	}
	configPaths = []string{
		".",
	}
)

// --- Configuration --- //

// Configuration struct
type Configuration struct {
	Debug          bool   `mapstructure:"DEBUG"`
	Port           int    `mapstructure:"PORT"`
	LogFile        string `mapstructure:"LOG_FILE"`
	Protocol       string
	ServerName     string      `mapstructure:"SERVER_NAME"`
	Auth           string      `mapstructure:"AUTH"`
	Client         string      `mapstructure:"CLIENT"`
	Endpoints      Endpoints   `mapstructure:",squash"`
	SSLCert        string      `mapstructure:"SSL_CERT"`
	SSLKey         string      `mapstructure:"SSL_KEY"`
	Db             DataSource  `mapstructure:",squash"`
	JWTKey         string      `mapstructure:"JWT_KEY"`
	RSAPublicKey   string      `mapstructure:"RSA_PUBLIC_KEY"`
	RSAPrivateKey  string      `mapstructure:"RSA_PRIVATE_KEY"`
	Redis          RedisConfig `mapstructure:",squash"`
	AllowedOrigins string      `mapstructure:"ALLOWED_ORIGINS"`
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

type RedisConfig struct {
	Host            string `mapstructure:"REDIS_HOST"`
	Port            int    `mapstructure:"REDIS_PORT"`
	Password        string `mapstructure:"REDIS_PASSWORD"`
	Db              int    `mapstructure:"REDIS_DB"`
	RedisExpSeconds int    `mapstructure:"REDIS_EXP_SECONDS"`
}

func ReadConfig(ENV string) (Configuration, error) {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	if ENV != "" {
		log.Printf("Running in ENV: %s", ENV)
		viper.SetConfigName(ENV)
		for _, p := range configPaths {
			viper.AddConfigPath(p)
		}
		err := viper.ReadInConfig()
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println("No ENV specified. Falling back to environment variables and defaults.")
	}
	viper.AutomaticEnv()
	var config Configuration
	err := viper.Unmarshal(&config)
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
