package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/cheebz/go-auth-helpers"
	"github.com/cheebz/go-pub/pkg/activitypub"
	"github.com/cheebz/go-pub/pkg/responses"
	"github.com/cheebz/go-pub/pkg/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ActivityPubMiddleware struct {
	auth             string
	response         responses.Response
	acceptMiddleware func(h http.Handler) http.Handler
}

func NewActivityPubMiddleware(_client string, _auth string, _response responses.Response) Middleware {
	m := &ActivityPubMiddleware{
		auth:             _auth,
		response:         _response,
		acceptMiddleware: createAcceptMiddleware(_client),
	}

	return m
}

func (m *ActivityPubMiddleware) CreateCORSMiddleware(allowedOrigins []string) func(h http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
	})
	return cors.Handler
}

func createAcceptMiddleware(client string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		if client != "" {
			serveClient := createServeClient(client)
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// check Accept header and serve client if it's not there
				err := activitypub.CheckAccept(r.Header)
				if err != nil {
					serveClient(w, r)
					return
				}
				h.ServeHTTP(w, r)
			})
		}
		return http.HandlerFunc(h.ServeHTTP)
	}
}

func createServeClient(client string) func(w http.ResponseWriter, r *http.Request) {
	if isValidURL(client) {
		return proxyToClient(client)
	}
	return serveStaticSite(client)
}

// isValidURL tests a string to determine if it is a well-structured url or not.
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func proxyToClient(client string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		origin, _ := url.Parse(client)
		director := func(req *http.Request) {
			req.URL.Scheme = origin.Scheme
			req.URL.Host = origin.Host
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	}
}

func serveStaticSite(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileRegexp := regexp.MustCompile(`\.[a-zA-Z]*$`)
		if !fileRegexp.MatchString(r.URL.Path) {
			// if not file, serve client app
			http.ServeFile(w, r, fmt.Sprintf("%s/index.html", path))
		} else {
			// else serve static file
			http.FileServer(http.Dir(path)).ServeHTTP(w, r)
		}
	}
}

func (m *ActivityPubMiddleware) AcceptMiddleware(h http.Handler) http.Handler {
	return m.acceptMiddleware(h)
}

func (m *ActivityPubMiddleware) ContentTypeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := activitypub.CheckContentType(r.Header)
		if err != nil {
			m.response.BadRequest(w, err)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (m *ActivityPubMiddleware) CreateJwtUsernameMiddleware(nameParam string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authMap, err := auth.Authenticate(w, r, m.auth)
			if err != nil {
				m.response.UnauthorizedRequest(w, err)
				return
			}
			if s, ok := authMap["username"].(string); !ok {
				m.response.UnauthorizedRequest(w, errors.New("invalid response from auth endpoint"))
				return
			} else {
				name := mux.Vars(r)[nameParam]
				if s != name {
					m.response.UnauthorizedRequest(w, errors.New("that's not yours"))
					return
				}
				h.ServeHTTP(w, r)
			}
		})
	}
}

func (m *ActivityPubMiddleware) CreateUserMiddleware(service services.Service) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authMap, err := auth.Authenticate(w, r, m.auth)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}
			if username, ok := authMap["username"].(string); ok {
				err = service.CheckUser(username)
				if err != nil {
					_, err = service.CreateUser(username)
					if err != nil {
						m.response.BadRequest(w, err)
						return
					}
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
