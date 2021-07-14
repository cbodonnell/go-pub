package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// TODO: Create a package for this middleware accepting the jwt key as a parameter

func checkJWTClaims(r *http.Request) (*JWTClaims, error) {
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
		return nil, err
	}
	return claims, nil
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

func refreshMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := checkJWTClaims(r)
		if err != nil {
			refresh(w, r)
		}
		h.ServeHTTP(w, r)
	})
}

func refresh(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	authReq, err := http.NewRequest("GET", fmt.Sprintf("%s/", config.Auth), nil)
	if err != nil {
		internalServerError(w, err)
		return
	}
	for _, cookie := range r.Cookies() {
		authReq.AddCookie(cookie)
	}
	authResp, err := client.Do(authReq)
	if err != nil {
		internalServerError(w, err)
		return
	}
	for _, cookie := range authResp.Cookies() {
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
	}
}
