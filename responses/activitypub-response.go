package responses

import (
	"net/http"

	"github.com/cheebz/go-pub/logging"
)

type ActivityPubResponse struct {
	debug bool
}

func NewActivityPubResponse(_debug bool) Response {
	return &ActivityPubResponse{
		debug: _debug,
	}
}

func (a *ActivityPubResponse) Created(w http.ResponseWriter, iri string) {
	w.Header().Add("Location", iri)
	w.WriteHeader(http.StatusCreated)
}

func (a *ActivityPubResponse) Accepted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}

func (a *ActivityPubResponse) BadRequest(w http.ResponseWriter, err error) {
	logging.LogCaller(err)
	var msg string
	if a.debug {
		msg = err.Error()
	} else {
		msg = "Bad request"
	}
	http.Error(w, msg, http.StatusBadRequest)
}

func (a *ActivityPubResponse) NotFound(w http.ResponseWriter, err error) {
	logging.LogCaller(err)
	var msg string
	if a.debug {
		msg = err.Error()
	} else {
		msg = "Not found"
	}
	http.Error(w, msg, http.StatusNotFound)
}

func (a *ActivityPubResponse) UnauthorizedRequest(w http.ResponseWriter, err error) {
	logging.LogCaller(err)
	var msg string
	if a.debug {
		msg = err.Error()
	} else {
		msg = "Unauthorized"
	}
	http.Error(w, msg, http.StatusUnauthorized)
}

func (a *ActivityPubResponse) InternalServerError(w http.ResponseWriter, err error) {
	logging.LogCaller(err)
	var msg string
	if a.debug {
		msg = err.Error()
	} else {
		msg = "Internal server error"
	}
	http.Error(w, msg, http.StatusInternalServerError)
}
