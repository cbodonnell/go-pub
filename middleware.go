package main

import (
	"regexp"
	"fmt"
	"net/http"
	"net/url"
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

func acceptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := checkAccept(r.Header)
		if err != nil {
			// if not requesting activity serve client app
			if isValidURL(config.Client) {
				// if url append request URI and redirect
				http.Redirect(w, r, config.Client+r.URL.RequestURI(), http.StatusSeeOther)
				return
			} else {
				// else try and serve static site
				fileRegexp := regexp.MustCompile(`\.[a-zA-Z]*$`)
				if !fileRegexp.MatchString(r.URL.Path) {
					// if not file, serve client app
					http.ServeFile(w, r, fmt.Sprintf("%s/index.html", config.Client))
				} else {
					// else serve static file
					http.FileServer(http.Dir(config.Client)).ServeHTTP(w, r)
				}
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func jwtMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := checkJWTClaims(r)
		if err != nil {
			refresh(w, r)
		}
		_, err = checkJWTClaims(r)
		if err != nil {
			unauthorizedRequest(w, err)
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

func userMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := checkJWTClaims(r)
		if err != nil {
			unauthorizedRequest(w, err)
			return
		}
		err = checkUser(claims.Username)
		if err != nil {
			_, err = createUser(claims.Username)
			if err != nil {
				badRequest(w, err)
			}
		}
		h.ServeHTTP(w, r)
	})
}
