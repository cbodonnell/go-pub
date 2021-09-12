package jwt

import "net/http"

type JWT interface {
	CheckJWTClaims(r *http.Request) (*JWTClaims, error)
	Refresh(w http.ResponseWriter, r *http.Request)
}
