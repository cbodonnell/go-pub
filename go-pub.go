package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cheebz/go-pub/cache"
	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/handlers"
	"github.com/cheebz/go-pub/logging"
	"github.com/cheebz/go-pub/middleware"
	"github.com/cheebz/go-pub/repositories"
	"github.com/cheebz/go-pub/resources"
	"github.com/cheebz/go-pub/responses"
	"github.com/cheebz/go-pub/services"
	"github.com/cheebz/go-pub/workers"
)

func main() {
	// Get configuration
	ENV := os.Getenv("ENV")
	conf, err := config.ReadConfig(ENV)
	if err != nil {
		log.Fatal(err)
	}

	// create cache layer
	cache := cache.NewRedisCache(conf)
	err = cache.FlushDB()
	if err != nil {
		log.Println(err)
	}
	// create repository
	repo := repositories.NewPSQLRepository(conf, cache)
	defer repo.Close()
	// create federation worker
	worker := workers.NewFederationWorker(conf, repo)
	go worker.Start()
	// create service
	service := services.NewActivityPubService(conf, repo, worker)
	// create response writer
	response := responses.NewActivityPubResponse(conf.Debug)
	// create middleware helper
	middle := middleware.NewActivityPubMiddleware(conf.Client, conf.Auth, response)
	// create resource generator
	resource := resources.NewActivityPubResource(conf)
	// create handler (TODO: Make an options struct??)
	handler := handlers.NewMuxHandler(conf.Endpoints, middle, service, resource, response)
	if conf.AllowedOrigins != "" {
		handler.AllowCORS(strings.Split(conf.AllowedOrigins, ","))
	}
	r := handler.GetRouter()

	// Run server
	log.Println(fmt.Sprintf("Serving on port %d", conf.Port))

	// Set log file
	if conf.LogFile != "" {
		logFile := logging.SetLogFile(conf.LogFile)
		defer logFile.Close()
	}

	// TLS
	if conf.SSLCert == "" {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), r))
	}
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", conf.Port), conf.SSLCert, conf.SSLKey, r))
}
