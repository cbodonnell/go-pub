package controllers

import "net/http"

type UserController interface {
	GetUser(w http.ResponseWriter, r *http.Request)
}
