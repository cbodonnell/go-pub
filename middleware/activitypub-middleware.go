package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/cheebz/go-pub/activitypub"
	"github.com/cheebz/go-pub/jwt"
	"github.com/cheebz/go-pub/responses"
	"github.com/cheebz/go-pub/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ActivityPubMiddleware struct {
	client   string
	response responses.Response
	jwt      jwt.JWT
}

func NewActivityPubMiddleware(_client string, _response responses.Response, _jwt jwt.JWT) Middleware {
	return &ActivityPubMiddleware{
		client:   _client,
		response: _response,
		jwt:      _jwt,
	}
}

func (m *ActivityPubMiddleware) CreateCORSMiddleware(allowedOrigins []string) func(h http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
	})
	return cors.Handler
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

func (m *ActivityPubMiddleware) AcceptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := activitypub.CheckAccept(r.Header)
		if err != nil {
			// if not requesting activity serve client app
			if isValidURL(m.client) {
				origin, _ := url.Parse(m.client)
				director := func(req *http.Request) {
					req.URL.Scheme = origin.Scheme
					req.URL.Host = origin.Host
				}
				proxy := &httputil.ReverseProxy{Director: director}
				proxy.ServeHTTP(w, r)
				return
			} else {
				// else try and serve static site
				fileRegexp := regexp.MustCompile(`\.[a-zA-Z]*$`)
				if !fileRegexp.MatchString(r.URL.Path) {
					// if not file, serve client app
					http.ServeFile(w, r, fmt.Sprintf("%s/index.html", m.client))
				} else {
					// else serve static file
					http.FileServer(http.Dir(m.client)).ServeHTTP(w, r)
				}
				return
			}
		}
		h.ServeHTTP(w, r)
	})
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

// func (m *ActivityPubMiddleware) JwtMiddleware(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		_, err := m.jwt.CheckJWTClaims(r)
// 		if err != nil {
// 			m.jwt.Refresh(w, r)
// 			_, err = m.jwt.CheckJWTClaims(r)
// 			if err != nil {
// 				m.response.UnauthorizedRequest(w, err)
// 				return
// 			}
// 		}
// 		h.ServeHTTP(w, r)
// 	})
// }

func (m *ActivityPubMiddleware) CreateJwtUsernameMiddleware(nameParam string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := m.jwt.CheckJWTClaims(r)
			if err != nil {
				m.jwt.Refresh(w, r)
				claims, err = m.jwt.CheckJWTClaims(r)
				if err != nil {
					m.response.UnauthorizedRequest(w, err)
					return
				}
			}
			// TODO: Get rid of the mux dependency here??
			name := mux.Vars(r)[nameParam]
			if claims.Username != name {
				m.response.UnauthorizedRequest(w, errors.New("that's not yours"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func (m *ActivityPubMiddleware) CreateUserMiddleware(service services.Service) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := m.jwt.CheckJWTClaims(r)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}
			err = service.CheckUser(claims.Username)
			if err != nil {
				_, err = service.CreateUser(claims.Username)
				if err != nil {
					m.response.BadRequest(w, err)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
