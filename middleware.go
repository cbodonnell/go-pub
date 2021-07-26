package main

import (
	"net/http"
)

func acceptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := checkAccept(r.Header)
		if err != nil {
			// serve client app
			// fmt.Println(r.URL.RequestURI())
			http.Redirect(w, r, config.Client+r.URL.RequestURI(), http.StatusSeeOther)
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
