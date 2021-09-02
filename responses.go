package main

import (
	"net/http"
)

func created(w http.ResponseWriter, iri string) {
	w.Header().Add("Location", iri)
	w.WriteHeader(http.StatusCreated)
}

func accepted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}

func badRequest(w http.ResponseWriter, err error) {
	logCaller(err)
	var msg string
	if config.Debug {
		msg = err.Error()
	} else {
		msg = "Bad request"
	}
	http.Error(w, msg, http.StatusBadRequest)
}

func notFound(w http.ResponseWriter, err error) {
	logCaller(err)
	var msg string
	if config.Debug {
		msg = err.Error()
	} else {
		msg = "Not found"
	}
	http.Error(w, msg, http.StatusNotFound)
}

func unauthorizedRequest(w http.ResponseWriter, err error) {
	logCaller(err)
	var msg string
	if config.Debug {
		msg = err.Error()
	} else {
		msg = "Unauthorized"
	}
	http.Error(w, msg, http.StatusUnauthorized)
}

func internalServerError(w http.ResponseWriter, err error) {
	logCaller(err)
	var msg string
	if config.Debug {
		msg = err.Error()
	} else {
		msg = "Internal server error"
	}
	http.Error(w, msg, http.StatusInternalServerError)
}
