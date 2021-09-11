package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/handlers"
	"github.com/cheebz/go-pub/logging"
	"github.com/cheebz/go-pub/repositories"
	"github.com/rs/cors"
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
	repo = repositories.NewPSQLRepository(conf.Db)
	defer repo.Close()
	// TODO: create services package
	// create handler
	r := handlers.NewMuxHandler(conf, repo)

	// TODO: Move remaining routes to mux-handlers package
	wf := r.NewRoute().Subrouter() // -> webfinger
	wf.HandleFunc("/.well-known/webfinger", getWebFinger).Methods("GET", "OPTIONS")

	get := r.NewRoute().Subrouter() // -> public GET requests
	get.Use(acceptMiddleware, userMiddleware)
	// get.HandleFunc("/users/{name:[[:alnum:]]+}", controller.GetUser).Methods("GET", "OPTIONS")
	get.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", getOutbox).Methods("GET", "OPTIONS")
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

	// Run server
	port := conf.Port
	log.Println(fmt.Sprintf("Serving on port %d", port))

	// CORS in dev
	if ENV == "dev" {
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			AllowCredentials: true,
		})
		r.Use(cors.Handler)
	}

	if conf.LogFile != "" {
		logFile := logging.SetLogFile(conf.LogFile)
		defer logFile.Close()
	}

	// TLS
	if conf.SSLCert == "" {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
	}
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), conf.SSLCert, conf.SSLKey, r))
}
