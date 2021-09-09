package responses

import (
	"net/http"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/logging"
)

var (
	debug bool = false
)

func Debug() {
	debug = true
}

func Created(w http.ResponseWriter, iri string) {
	w.Header().Add("Location", iri)
	w.WriteHeader(http.StatusCreated)
}

func Accepted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}

func BadRequest(w http.ResponseWriter, err error) {
	logging.LogCaller(err)
	var msg string
	if config.C.Debug {
		msg = err.Error()
	} else {
		msg = "Bad request"
	}
	http.Error(w, msg, http.StatusBadRequest)
}

func NotFound(w http.ResponseWriter, err error) {
	logging.LogCaller(err)
	var msg string
	if config.C.Debug {
		msg = err.Error()
	} else {
		msg = "Not found"
	}
	http.Error(w, msg, http.StatusNotFound)
}

func UnauthorizedRequest(w http.ResponseWriter, err error) {
	logging.LogCaller(err)
	var msg string
	if config.C.Debug {
		msg = err.Error()
	} else {
		msg = "Unauthorized"
	}
	http.Error(w, msg, http.StatusUnauthorized)
}

func InternalServerError(w http.ResponseWriter, err error) {
	logging.LogCaller(err)
	var msg string
	if config.C.Debug {
		msg = err.Error()
	} else {
		msg = "Internal server error"
	}
	http.Error(w, msg, http.StatusInternalServerError)
}
