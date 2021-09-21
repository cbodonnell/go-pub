package handlers

import (
	"net/http"
)

type Handler interface {
	GetRouter() http.Handler
	AllowCORS(allowedOrigins []string)
	GetWebFinger(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetInbox(w http.ResponseWriter, r *http.Request)
	GetOutbox(w http.ResponseWriter, r *http.Request)
	GetFollowing(w http.ResponseWriter, r *http.Request)
	GetFollowers(w http.ResponseWriter, r *http.Request)
	GetLiked(w http.ResponseWriter, r *http.Request)
	GetActivity(w http.ResponseWriter, r *http.Request)
	GetObject(w http.ResponseWriter, r *http.Request)
	PostInbox(w http.ResponseWriter, r *http.Request)
	PostOutbox(w http.ResponseWriter, r *http.Request)
	SinkHandler(w http.ResponseWriter, r *http.Request)
}
