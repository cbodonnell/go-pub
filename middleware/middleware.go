package middleware

import (
	"net/http"

	"github.com/cheebz/go-pub/services"
)

type Middleware interface {
	CreateCORSMiddleware(allowedOrigins []string) func(h http.Handler) http.Handler
	AcceptMiddleware(h http.Handler) http.Handler
	ContentTypeMiddleware(h http.Handler) http.Handler
	// JwtMiddleware(h http.Handler) http.Handler
	CreateUserMiddleware(service services.Service) func(h http.Handler) http.Handler
	CreateJwtUsernameMiddleware(name string) func(h http.Handler) http.Handler
}
