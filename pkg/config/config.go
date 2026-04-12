package config

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var (
	defaults = map[string]interface{}{
		"DEBUG":                 true,
		"PORT":                  8080,
		"LOG_FILE":              "",
		"SERVER_NAME":           "localhost",
		"PROTOCOL":              "http",
		"CLIENT":                "",
		"ENDPOINT_USERS":        "users",
		"ENDPOINT_ACTIVITIES":   "activities",
		"ENDPOINT_OBJECTS":      "objects",
		"ENDPOINT_FEED":         "feed",
		"ENDPOINT_INBOX":        "inbox",
		"ENDPOINT_OUTBOX":       "outbox",
		"ENDPOINT_FOLLOWING":    "following",
		"ENDPOINT_FOLLOWERS":    "followers",
		"ENDPOINT_LIKED":        "liked",
		"ENDPOINT_UPLOAD_MEDIA": "uploadMedia",
		"ENDPOINT_UPLOADS":      "uploads",
		"ENDPOINT_LINKS":        "links",
		"ENDPOINT_CHECK":        "check",
		"UPLOAD_DIR":            "./uploads/",
		"SSL_CERT":              "",
		"SSL_KEY":               "",
		"DB_HOST":               "localhost",
		"DB_PORT":               5432,
		"DB_NAME":               "gopub",
		"DB_USER":               "gopub",
		"DB_PASSWORD":           "gopub",
		"JWT_KEY":               "changeme",
		"JWT_MAX_AGE":           1200,
		"REFRESH_MAX_AGE":       2592000,
		"RSA_PUBLIC_KEY":        "./keys/public.pem",
		"RSA_PRIVATE_KEY":       "./keys/private.pem",
		"REDIS_HOST":            "localhost",
		"REDIS_PORT":            6379,
		"REDIS_PASSWORD":        "",
		"REDIS_DB":              0,
		"REDIS_EXP_SECONDS":     3600,
		"ALLOWED_ORIGINS":       "",
		"PAGE_LENGTH":           10,
		"REGISTER":              true,
	}
	configPaths = []string{
		".",
	}
)

// Configuration struct
type Configuration struct {
	Debug          bool        `mapstructure:"DEBUG"`
	Port           int         `mapstructure:"PORT"`
	LogFile        string      `mapstructure:"LOG_FILE"`
	Protocol       string      `mapstructure:"PROTOCOL"`
	ServerName     string      `mapstructure:"SERVER_NAME"`
	Client         string      `mapstructure:"CLIENT"`
	Endpoints      Endpoints   `mapstructure:",squash"`
	UploadDir      string      `mapstructure:"UPLOAD_DIR"`
	SSLCert        string      `mapstructure:"SSL_CERT"`
	SSLKey         string      `mapstructure:"SSL_KEY"`
	Db             DataSource  `mapstructure:",squash"`
	JWTKey         string      `mapstructure:"JWT_KEY"`
	JWTMaxAge      int         `mapstructure:"JWT_MAX_AGE"`
	RefreshMaxAge  int         `mapstructure:"REFRESH_MAX_AGE"`
	RSAPublicKey   string      `mapstructure:"RSA_PUBLIC_KEY"`
	RSAPrivateKey  string      `mapstructure:"RSA_PRIVATE_KEY"`
	Redis          RedisConfig `mapstructure:",squash"`
	AllowedOrigins string      `mapstructure:"ALLOWED_ORIGINS"`
	PageLength     int         `mapstructure:"PAGE_LENGTH"`
	Register       bool        `mapstructure:"REGISTER"`
}

type Endpoints struct {
	Users       string `mapstructure:"ENDPOINT_USERS"`
	Activities  string `mapstructure:"ENDPOINT_ACTIVITIES"`
	Objects     string `mapstructure:"ENDPOINT_OBJECTS"`
	Feed        string `mapstructure:"ENDPOINT_FEED"`
	Inbox       string `mapstructure:"ENDPOINT_INBOX"`
	Outbox      string `mapstructure:"ENDPOINT_OUTBOX"`
	Following   string `mapstructure:"ENDPOINT_FOLLOWING"`
	Followers   string `mapstructure:"ENDPOINT_FOLLOWERS"`
	Liked       string `mapstructure:"ENDPOINT_LIKED"`
	UploadMedia string `mapstructure:"ENDPOINT_UPLOAD_MEDIA"`
	Uploads     string `mapstructure:"ENDPOINT_UPLOADS"`
	Links       string `mapstructure:"ENDPOINT_LINKS"`
	Check       string `mapstructure:"ENDPOINT_CHECK"`
}

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

	if config.Protocol == "" {
		if config.SSLCert == "" {
			config.Protocol = "http"
		} else {
			config.Protocol = "https"
		}
	}

	// Auto-generate RSA keys if they don't exist
	if err := ensureRSAKeys(config.RSAPublicKey, config.RSAPrivateKey); err != nil {
		return config, fmt.Errorf("failed to ensure RSA keys: %w", err)
	}

	// Read RSA key contents
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

// ensureRSAKeys generates a 2048-bit RSA key pair at the given paths if the
// private key file does not already exist.
func ensureRSAKeys(publicKeyPath, privateKeyPath string) error {
	if _, err := os.Stat(privateKeyPath); err == nil {
		return nil // already exists
	}

	log.Println("RSA keys not found — generating new key pair...")

	if err := os.MkdirAll(filepath.Dir(privateKeyPath), 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(publicKeyPath), 0700); err != nil {
		return err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Write private key
	privFile, err := os.OpenFile(privateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer privFile.Close()
	if err := pem.Encode(privFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return err
	}

	// Write public key
	pubFile, err := os.OpenFile(publicKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer pubFile.Close()
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	if err := pem.Encode(pubFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}); err != nil {
		return err
	}

	log.Printf("Generated RSA key pair: %s, %s", privateKeyPath, publicKeyPath)
	return nil
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
