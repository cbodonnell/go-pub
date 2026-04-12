package middleware

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/cheebz/go-pub/pkg/activitypub"
	"github.com/cheebz/go-pub/pkg/auth"
	"github.com/cheebz/go-pub/pkg/responses"
	"github.com/cheebz/go-pub/pkg/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ActivityPubMiddleware struct {
	authService      auth.Service
	response         responses.Response
	acceptMiddleware func(h http.Handler) http.Handler
}

// NewActivityPubMiddleware creates the middleware.
// client can be:
//   - a URL ("http://...") → reverse proxy to that address
//   - a directory path ("./client/dist") → serve static files
//   - empty string → fall back to embeddedFS (nil embeddedFS = 404 for non-AP requests)
func NewActivityPubMiddleware(client string, embeddedFS fs.FS, authService auth.Service, response responses.Response) Middleware {
	m := &ActivityPubMiddleware{
		authService:      authService,
		response:         response,
		acceptMiddleware: createAcceptMiddleware(client, embeddedFS),
	}
	return m
}

func (m *ActivityPubMiddleware) CreateCORSMiddleware(allowedOrigins []string) func(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
	})
	return c.Handler
}

func createAcceptMiddleware(client string, embeddedFS fs.FS) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		serveClient := createServeClient(client, embeddedFS)
		if serveClient == nil {
			return http.HandlerFunc(h.ServeHTTP)
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := activitypub.CheckAccept(r.Header); err != nil {
				serveClient(w, r)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func createServeClient(client string, embeddedFS fs.FS) func(w http.ResponseWriter, r *http.Request) {
	if client != "" {
		if isValidURL(client) {
			return proxyToClient(client)
		}
		return serveStaticSite(client)
	}
	if embeddedFS != nil {
		return serveEmbedded(embeddedFS)
	}
	return nil
}

func isValidURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func proxyToClient(client string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		origin, _ := url.Parse(client)
		proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
			req.URL.Scheme = origin.Scheme
			req.URL.Host = origin.Host
		}}
		proxy.ServeHTTP(w, r)
	}
}

func serveStaticSite(path string) func(w http.ResponseWriter, r *http.Request) {
	fileRegexp := regexp.MustCompile(`\.[a-zA-Z]*$`)
	return func(w http.ResponseWriter, r *http.Request) {
		if !fileRegexp.MatchString(r.URL.Path) {
			http.ServeFile(w, r, fmt.Sprintf("%s/index.html", path))
		} else {
			http.FileServer(http.Dir(path)).ServeHTTP(w, r)
		}
	}
}

func serveEmbedded(fsys fs.FS) func(w http.ResponseWriter, r *http.Request) {
	fileServer := http.FileServer(http.FS(fsys))
	fileRegexp := regexp.MustCompile(`\.[a-zA-Z0-9]+$`)
	return func(w http.ResponseWriter, r *http.Request) {
		if !fileRegexp.MatchString(r.URL.Path) {
			// SPA: all non-asset paths serve index.html so client-side routing works
			data, err := fs.ReadFile(fsys, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data) //nolint:errcheck
			return
		}
		fileServer.ServeHTTP(w, r)
	}
}

func (m *ActivityPubMiddleware) AcceptMiddleware(h http.Handler) http.Handler {
	return m.acceptMiddleware(h)
}

func (m *ActivityPubMiddleware) ContentTypeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := activitypub.CheckContentType(r.Header); err != nil {
			m.response.BadRequest(w, err)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// CreateJwtUsernameMiddleware guards routes that require the authenticated user
// to match the {name} path parameter.
func (m *ActivityPubMiddleware) CreateJwtUsernameMiddleware(nameParam string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := m.authService.Authenticate(w, r)
			if err != nil {
				m.response.UnauthorizedRequest(w, errors.New("authentication required"))
				return
			}
			name := mux.Vars(r)[nameParam]
			if claims.Username != name {
				m.response.UnauthorizedRequest(w, errors.New("that's not yours"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

// CreateUserMiddleware auto-creates an ActivityPub user record for any
// authenticated request where one doesn't exist yet.
func (m *ActivityPubMiddleware) CreateUserMiddleware(service services.Service) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := m.authService.GetClaims(r)
			if err == nil {
				if err := service.CheckUser(claims.Username); err != nil {
					if _, err := service.CreateUser(claims.Username); err != nil {
						m.response.BadRequest(w, err)
						return
					}
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
