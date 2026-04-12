package commands

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cheebz/go-pub/internal/static"
	"github.com/cheebz/go-pub/pkg/activitypub"
	"github.com/cheebz/go-pub/pkg/auth"
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

	// Get configuration (also auto-generates RSA keys if absent)
	ENV := os.Getenv("ENV")
	conf, err := config.ReadConfig(ENV)
	if err != nil {
		return err
	}

	// Auth
	authRepo := auth.NewPSQLRepository(conf)
	jwtHelper := auth.NewJWTHelper(conf.JWTKey, conf.JWTMaxAge, conf.RefreshMaxAge, conf.SSLCert != "")
	authService := auth.NewService(authRepo, jwtHelper, conf.RefreshMaxAge)
	authHandler := auth.NewHandler(authService, conf)
	go auth.NewPurgeWorker(authRepo).Start()

	// Cache
	redisCache := cache.NewRedisCache(conf)
	if err := redisCache.FlushDB(); err != nil {
		log.Println("failed to flush cache:", err)
	}

	// ActivityPub
	repo := repositories.NewPSQLRepository(conf, redisCache)
	defer repo.Close()
	fileWorker := workers.NewFileWorker(conf, repo)
	go fileWorker.Start()
	federator := activitypub.NewFederator(conf, repo)
	service := services.NewActivityPubService(conf, repo, federator)
	response := responses.NewActivityPubResponse(conf.Debug)
	middle := middleware.NewActivityPubMiddleware(conf.Client, static.FS, authService, response)
	resource := resources.NewActivityPubResource(conf)
	handler := handlers.NewMuxHandler(conf, middle, service, resource, response, authHandler)
	if conf.AllowedOrigins != "" {
		handler.AllowCORS(strings.Split(conf.AllowedOrigins, ","))
	}
	r := handler.GetRouter()

	log.Printf("Serving on port %d\n", conf.Port)

	if conf.LogFile != "" {
		logFile := logging.SetLogFile(conf.LogFile)
		defer logFile.Close()
	}

	if conf.SSLCert == "" {
		return http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), r)
	}
	return http.ListenAndServeTLS(fmt.Sprintf(":%d", conf.Port), conf.SSLCert, conf.SSLKey, r)
}
