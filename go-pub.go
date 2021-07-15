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

	// Static home page for testing auth
	h := r.HandleFunc("/", home).Methods("GET").Subrouter()
	h.Use(refreshMiddleware)

	// This is a client-to-server GET of an activity
	g := r.Methods("GET").Subrouter()
	g.HandleFunc("/.well-known/webfinger", getWebFinger).Methods("GET")
	g.HandleFunc("/users/{name:[[:alnum:]]+}", getUser).Methods("GET")
	g.HandleFunc("/users/{name:[[:alnum:]]+}/inbox", getInbox).Methods("GET")
	g.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", getOutbox).Methods("GET")
	g.HandleFunc("/users/{name:[[:alnum:]]+}/following", getFollowing).Methods("GET")
	g.HandleFunc("/users/{name:[[:alnum:]]+}/followers", getFollowers).Methods("GET")
	g.HandleFunc("/users/{name:[[:alnum:]]+}/liked", getLiked).Methods("GET")

	// This is a client-to-server POST of an activity
	p := r.Methods("POST").Subrouter()
	p.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", postOutbox).Methods("POST")
	// p.Use(jwtMiddleware)

	// Static files
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("./static/"))))
	// r.Use(jwtMiddleware)

	// Run server
	port := config.Port
	fmt.Println(fmt.Sprintf("Serving on port %d", port))

	// CORS in dev
	if ENV == "dev" {
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:4200", "http://127.0.0.1:4200"},
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
