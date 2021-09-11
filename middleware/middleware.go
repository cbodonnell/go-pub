package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/cheebz/go-pub/activitypub"
	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/jwt"
	"github.com/cheebz/go-pub/repositories"
	"github.com/cheebz/go-pub/responses"
)

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

func AcceptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := activitypub.CheckAccept(r.Header)
		if err != nil {
			// if not requesting activity serve client app
			if isValidURL(config.C.Client) {
				// if url append request URI and redirect
				http.Redirect(w, r, config.C.Client+r.URL.RequestURI(), http.StatusSeeOther)
				return
			} else {
				// else try and serve static site
				fileRegexp := regexp.MustCompile(`\.[a-zA-Z]*$`)
				if !fileRegexp.MatchString(r.URL.Path) {
					// if not file, serve client app
					http.ServeFile(w, r, fmt.Sprintf("%s/index.html", config.C.Client))
				} else {
					// else serve static file
					http.FileServer(http.Dir(config.C.Client)).ServeHTTP(w, r)
				}
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func ContentTypeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := activitypub.CheckContentType(r.Header)
		if err != nil {
			responses.BadRequest(w, err)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func JwtMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := jwt.CheckJWTClaims(r)
		if err != nil {
			jwt.Refresh(w, r)
		}
		_, err = jwt.CheckJWTClaims(r)
		if err != nil {
			responses.UnauthorizedRequest(w, err)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// func refreshMiddleware(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		_, err := checkJWTClaims(r)
// 		if err != nil {
// 			refresh(w, r)
// 		}
// 		h.ServeHTTP(w, r)
// 	})
// }

func CreateUserMiddleware(repo repositories.Repository) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := jwt.CheckJWTClaims(r)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}
			err = repo.CheckUser(claims.Username)
			if err != nil {
				_, err = repo.CreateUser(claims.Username)
				if err != nil {
					responses.BadRequest(w, err)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

// func UserMiddleware(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		claims, err := jwt.CheckJWTClaims(r)
// 		if err != nil {
// 			h.ServeHTTP(w, r)
// 			return
// 		}
// 		err = service.checkUser(claims.Username)
// 		if err != nil {
// 			_, err = service.createUser(claims.Username)
// 			if err != nil {
// 				responses.BadRequest(w, err)
// 				return
// 			}
// 		}
// 		h.ServeHTTP(w, r)
// 	})
// }
