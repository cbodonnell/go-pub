package main

import (
	"fmt"
	"net/http"

	"github.com/cheebz/go-pub/config"
	"github.com/dgrijalva/jwt-go"
)

// TODO: Create a package for this middleware accepting the jwt key as a parameter

// JWTClaims struct
type JWTClaims struct {
	UserID   int     `json:"user_id"`
	Username string  `json:"username"`
	UUID     string  `json:"uuid"`
	Groups   []Group `json:"groups"`
	jwt.StandardClaims
}

func checkJWTClaims(r *http.Request) (*JWTClaims, error) {
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		return nil, err
	}
	tokenString := jwtCookie.Value

	claims := &JWTClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.C.JWTKey), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func refresh(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	authReq, err := http.NewRequest("GET", fmt.Sprintf("%s/", config.C.Auth), nil)
	if err != nil {
		return
	}
	for _, cookie := range r.Cookies() {
		authReq.AddCookie(cookie)
	}
	authResp, err := client.Do(authReq)
	if err != nil {
		return
	}
	for _, cookie := range authResp.Cookies() {
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
	}
}
