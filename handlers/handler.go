package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	GetRouter() *mux.Router // TODO: Make more generic once able to
	AllowCORS(allowedOrigins []string)
	GetWebFinger(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetOutbox(w http.ResponseWriter, r *http.Request)
}
