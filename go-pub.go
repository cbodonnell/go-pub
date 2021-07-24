package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Get configuration
	ENV := os.Getenv("ENV")
	if ENV == "" {
		ENV = "dev"
	}
	fmt.Println(fmt.Sprintf("Running in ENV: %s", ENV))
	config = getConfig(ENV)

	db = connectDb(config.Db)
	defer db.Close()

	// Init router
	r := mux.NewRouter()
	pub := r.NewRoute().Subrouter()
	auth := r.NewRoute().Subrouter()

	// TODO: Break subrouters into public and auth?
	// public -> GET from Outbox and POST to Inbox
	// auth -> POST to Outbox and GET from Inbox

	// This is a client-to-server GET of an activity
	// g := r.Methods("GET").Subrouter()
	pub.HandleFunc("/", home).Methods("GET")
	pub.HandleFunc("/register", register).Methods("GET")

	pub.HandleFunc("/.well-known/webfinger", getWebFinger).Methods("GET")
	pub.HandleFunc("/users/{name:[[:alnum:]]+}", getUser).Methods("GET")
	pub.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", getOutbox).Methods("GET")
	pub.HandleFunc("/users/{name:[[:alnum:]]+}/following", getFollowing).Methods("GET")
	pub.HandleFunc("/users/{name:[[:alnum:]]+}/followers", getFollowers).Methods("GET")
	pub.HandleFunc("/users/{name:[[:alnum:]]+}/liked", getLiked).Methods("GET")
	// g.Use(refreshMiddleware)

	// This is a client-to-server POST of an activity
	// p := r.Methods("POST", "OPTIONS").Subrouter()
	auth.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", postOutbox).Methods("POST", "OPTIONS")
	auth.HandleFunc("/users/{name:[[:alnum:]]+}/inbox", getInbox).Methods("GET")
	auth.Use(jwtMiddleware, userMiddleware)

	// Static files
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("./static/"))))
	// r.Use(jwtMiddleware)

	// Run server
	port := config.Port
	fmt.Println(fmt.Sprintf("Serving on port %d", port))

	// CORS in dev
	if ENV == "dev" {
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			AllowCredentials: true,
		})
		r.Use(cors.Handler)
	}

	// TLS
	if config.SSLCert == "" {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
	}
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), config.SSLCert, config.SSLKey, r))
}
