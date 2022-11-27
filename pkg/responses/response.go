package responses

import "net/http"

type Response interface {
	Created(w http.ResponseWriter, iri string)
	Accepted(w http.ResponseWriter)
	BadRequest(w http.ResponseWriter, err error)
	NotFound(w http.ResponseWriter, err error)
	UnauthorizedRequest(w http.ResponseWriter, err error)
	InternalServerError(w http.ResponseWriter, err error)
}
