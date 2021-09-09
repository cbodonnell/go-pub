package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/controllers"
	"github.com/cheebz/go-pub/repositories"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var logChan = make(chan string)
var fedChan = make(chan Federation)

func main() {

	// get config
	// connect db
	// create repository
	// create service
	// create controller
	// create router
	// handle routes
	// apply middleware
	// listen and serve

	// TODO: Create config package
	// Get configuration
	ENV := os.Getenv("ENV")
	if ENV == "" {
		ENV = "dev"
	}
	log.Println(fmt.Sprintf("Running in ENV: %s", ENV))
	// config = getConfig(ENV)
	config.ReadConfig(ENV)

	db = repositories.ConnectDb(config.C.Db)
	defer db.Close()
	repository = repositories.NewPSQLRepository(db)
	controller := controllers.NewUserController(repository)

	// Init router
	r := mux.NewRouter()
	wf := r.NewRoute().Subrouter()   // -> webfinger
	get := r.NewRoute().Subrouter()  // -> public GET requests
	post := r.NewRoute().Subrouter() // -> public POST requests
	// auth := r.NewRoute().Subrouter() // -> POST to Outbox and GET from Inbox
	sink := r.NewRoute().Subrouter() // -> sink to handle all other routes

	wf.HandleFunc("/.well-known/webfinger", getWebFinger).Methods("GET", "OPTIONS")

	get.HandleFunc("/users/{name:[[:alnum:]]+}", controller.GetUser).Methods("GET", "OPTIONS")
	get.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", getOutbox).Methods("GET", "OPTIONS")
	get.HandleFunc("/users/{name:[[:alnum:]]+}/following", getFollowing).Methods("GET", "OPTIONS")
	get.HandleFunc("/users/{name:[[:alnum:]]+}/followers", getFollowers).Methods("GET", "OPTIONS")
	get.HandleFunc("/users/{name:[[:alnum:]]+}/liked", getLiked).Methods("GET", "OPTIONS")
	get.HandleFunc("/activities/{id}", getActivity).Methods("GET", "OPTIONS")
	get.HandleFunc("/objects/{id}", getObject).Methods("GET", "OPTIONS")
	get.Use(acceptMiddleware, userMiddleware)

	post.HandleFunc("/users/{name:[[:alnum:]]+}/inbox", postInbox).Methods("POST", "OPTIONS")
	post.Use(contentTypeMiddleware, userMiddleware)

	aGet := get.NewRoute().Subrouter()
	aGet.HandleFunc("/users/{name:[[:alnum:]]+}/inbox", getInbox).Methods("GET", "OPTIONS")
	aGet.Use(jwtMiddleware)

	aPost := post.NewRoute().Subrouter()
	aPost.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", postOutbox).Methods("POST", "OPTIONS")
	aPost.Use(jwtMiddleware)

	sink.PathPrefix("/").HandlerFunc(sinkHandler).Methods("GET", "OPTIONS")
	sink.Use(acceptMiddleware)

	// Static files
	// TODO: This should be done in a more sure way with permissions checking
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("./static/"))))
	// r.Use(jwtMiddleware)

	// TODO: Start federation worker
	go handleLogs()
	go handleFederation()

	// Run server
	port := config.C.Port
	log.Println(fmt.Sprintf("Serving on port %d", port))

	// CORS in dev
	if ENV == "dev" {
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			AllowCredentials: true,
		})
		r.Use(cors.Handler)
	}

	if config.C.LogFile != "" {
		logFile := setLogFile(config.C.LogFile)
		defer logFile.Close()
	}

	// TLS
	if config.C.SSLCert == "" {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
	}
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), config.C.SSLCert, config.C.SSLKey, r))
}
