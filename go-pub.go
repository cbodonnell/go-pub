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

	// Init router
	r := mux.NewRouter()
	r.HandleFunc("/.well-known/webfinger", webFinger)
	r.HandleFunc("/users/{name:[[:alnum:]]+}", getUser)
	r.HandleFunc("/users/{name:[[:alnum:]]+}/inbox", getInbox)
	r.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", getOutbox)
	// r.HandleFunc("/users/{name:[[:alnum:]]+}/following", getFollowing)
	// r.HandleFunc("/users/{name:[[:alnum:]]+}/followers", getFollowers)
	// r.HandleFunc("/users/{name:[[:alnum:]]+}/liked", getLiked)
	// r.Use(jwtMiddleware)

	// Run server
	port := config.Port
	fmt.Println(fmt.Sprintf("Serving on port %d", port))

	if ENV == "dev" {
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			AllowCredentials: true,
		})
		r.Use(cors.Handler)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
	}

	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), config.SSLCert, config.SSLKey, r))
}
