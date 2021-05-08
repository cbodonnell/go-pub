package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// TODO: Create a package for this middleware accepting the jwt key as a parameter

func checkClaims(r *http.Request) (*JWTClaims, error) {
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		return nil, err
	}
	tokenString := jwtCookie.Value

	claims := &JWTClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTKey), nil
	})
	if err != nil {
		return claims, err
	}
	return claims, nil
}

func jwtMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := checkClaims(r)
		if err != nil {
			unauthorizedRequest(w, err)
			return
		}
		h.ServeHTTP(w, r)
	})
}
