package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/handlers"
	"github.com/cheebz/go-pub/logging"
	"github.com/cheebz/go-pub/middleware"
	"github.com/cheebz/go-pub/repositories"
	"github.com/cheebz/go-pub/resources"
	"github.com/cheebz/go-pub/responses"
	"github.com/cheebz/go-pub/services"
)

var logChan = make(chan string)
var fedChan = make(chan Federation)

func main() {
	// Get configuration
	ENV := os.Getenv("ENV")
	if ENV == "" {
		ENV = "dev"
	}
	log.Println(fmt.Sprintf("Running in ENV: %s", ENV))
	conf, err := config.ReadConfig(ENV)
	if err != nil {
		log.Fatal("unable to read config")
	}

	// create repository
	repo := repositories.NewPSQLRepository(conf.Db)
	defer repo.Close()
	// create service
	service := services.NewActivityPubService(conf, repo)
	// create response writer
	response := responses.NewActivityPubResponse(conf.Debug)
	// create middleware engine
	middle := middleware.NewActivityPubMiddleware(conf.Client, response)
	// create resource generator
	resource := resources.NewActivityPubResource(conf)
	// create handler
	handler := handlers.NewMuxHandler(conf.Endpoints, middle, service, resource, response)
	if ENV == "dev" {
		handler.AllowCORS([]string{conf.Client})
	}
	r := handler.GetRouter()

	// TODO: Move remaining routes to mux-handlers package
	// wf := r.NewRoute().Subrouter() // -> webfinger
	// wf.HandleFunc("/.well-known/webfinger", getWebFinger).Methods("GET", "OPTIONS")

	get := r.NewRoute().Subrouter() // -> public GET requests
	get.Use(acceptMiddleware, userMiddleware)
	// get.HandleFunc("/users/{name:[[:alnum:]]+}", controller.GetUser).Methods("GET", "OPTIONS")
	// get.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", getOutbox).Methods("GET", "OPTIONS")
	get.HandleFunc("/users/{name:[[:alnum:]]+}/following", getFollowing).Methods("GET", "OPTIONS")
	get.HandleFunc("/users/{name:[[:alnum:]]+}/followers", getFollowers).Methods("GET", "OPTIONS")
	get.HandleFunc("/users/{name:[[:alnum:]]+}/liked", getLiked).Methods("GET", "OPTIONS")
	get.HandleFunc("/activities/{id}", getActivity).Methods("GET", "OPTIONS")
	get.HandleFunc("/objects/{id}", getObject).Methods("GET", "OPTIONS")

	post := r.NewRoute().Subrouter() // -> public POST requests
	post.Use(contentTypeMiddleware, userMiddleware)
	post.HandleFunc("/users/{name:[[:alnum:]]+}/inbox", postInbox).Methods("POST", "OPTIONS")

	aGet := get.NewRoute().Subrouter()
	aGet.Use(jwtMiddleware)
	aGet.HandleFunc("/users/{name:[[:alnum:]]+}/inbox", getInbox).Methods("GET", "OPTIONS")

	aPost := post.NewRoute().Subrouter()
	aPost.Use(jwtMiddleware)
	aPost.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", postOutbox).Methods("POST", "OPTIONS")

	sink := r.NewRoute().Subrouter() // -> sink to handle all other routes
	sink.Use(acceptMiddleware)
	sink.PathPrefix("/").HandlerFunc(sinkHandler).Methods("GET", "OPTIONS")

	// Static files
	// TODO: This should be done in a more sure way with permissions checking
	// r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("./static/"))))
	// r.Use(jwtMiddleware)

	go handleLogs()
	go handleFederation()

	if conf.LogFile != "" {
		logFile := logging.SetLogFile(conf.LogFile)
		defer logFile.Close()
	}

	// Run server
	log.Println(fmt.Sprintf("Serving on port %d", conf.Port))

	// TLS
	if conf.SSLCert == "" {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), r))
	}
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", conf.Port), conf.SSLCert, conf.SSLKey, r))
}
