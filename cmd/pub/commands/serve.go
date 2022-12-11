package commands

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cheebz/go-pub/pkg/activitypub"
	"github.com/cheebz/go-pub/pkg/cache"
	"github.com/cheebz/go-pub/pkg/config"
	"github.com/cheebz/go-pub/pkg/handlers"
	"github.com/cheebz/go-pub/pkg/logging"
	"github.com/cheebz/go-pub/pkg/middleware"
	"github.com/cheebz/go-pub/pkg/repositories"
	"github.com/cheebz/go-pub/pkg/resources"
	"github.com/cheebz/go-pub/pkg/responses"
	"github.com/cheebz/go-pub/pkg/services"
	"github.com/cheebz/go-pub/pkg/workers"
)

func Serve() error {
	serveCmd := flag.NewFlagSet(os.Args[1], flag.ExitOnError)

	serveCmd.Usage = func() {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintf(os.Stderr, "Usage: %s %s [flags]\n", os.Args[0], os.Args[1])
		fmt.Fprintln(os.Stderr, "")
		fmt.Println("Start the server")
		fmt.Fprintln(os.Stderr, "")
	}

	serveCmd.Parse(os.Args[2:])

	// Get configuration
	ENV := os.Getenv("ENV")
	conf, err := config.ReadConfig(ENV)
	if err != nil {
		return err
	}

	// create cache layer
	cache := cache.NewRedisCache(conf)
	err = cache.FlushDB()
	if err != nil {
		log.Println("failed to flush cache:", err)
	}
	// create repository
	repo := repositories.NewPSQLRepository(conf, cache)
	defer repo.Close()
	// create file worker
	fileWorker := workers.NewFileWorker(conf, repo)
	go fileWorker.Start()
	// create federator
	federator := activitypub.NewFederator(conf, repo)
	// create service
	service := services.NewActivityPubService(conf, repo, federator)
	// create response writer
	response := responses.NewActivityPubResponse(conf.Debug)
	// create middleware helper
	middle := middleware.NewActivityPubMiddleware(conf.Client, conf.Auth, response)
	// create resource generator
	resource := resources.NewActivityPubResource(conf)
	// create handler (TODO: Make an options struct??)
	handler := handlers.NewMuxHandler(conf, middle, service, resource, response)
	if conf.AllowedOrigins != "" {
		handler.AllowCORS(strings.Split(conf.AllowedOrigins, ","))
	}
	r := handler.GetRouter()

	// Run server
	log.Printf("Serving on port %d\n", conf.Port)

	// Set log file
	if conf.LogFile != "" {
		logFile := logging.SetLogFile(conf.LogFile)
		defer logFile.Close()
	}

	// TLS
	if conf.SSLCert == "" {
		return http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), r)
	}
	return http.ListenAndServeTLS(fmt.Sprintf(":%d", conf.Port), conf.SSLCert, conf.SSLKey, r)
}
