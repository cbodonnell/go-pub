package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cheebz/go-pub/activitypub"
	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/middleware"
	"github.com/cheebz/go-pub/resources"
	"github.com/cheebz/go-pub/responses"
	"github.com/cheebz/go-pub/services"
	"github.com/gorilla/mux"
)

type MuxHandler struct {
	endpoints  config.Endpoints
	middleware middleware.Middleware
	service    services.Service
	resource   resources.Resource
	response   responses.Response
	router     *mux.Router
}

func NewMuxHandler(_endpoints config.Endpoints, _middleware middleware.Middleware, _service services.Service, _resource resources.Resource, _response responses.Response) Handler {
	h := &MuxHandler{
		endpoints:  _endpoints,
		middleware: _middleware,
		service:    _service,
		resource:   _resource,
		response:   _response,
		router:     mux.NewRouter(),
	}

	wf := h.router.NewRoute().Subrouter() // -> webfinger
	wf.HandleFunc("/.well-known/webfinger", h.GetWebFinger).Methods("GET", "OPTIONS")

	userMiddleware := h.middleware.CreateUserMiddleware(h.service)

	get := h.router.NewRoute().Subrouter() // -> public GET requests
	get.Use(h.middleware.AcceptMiddleware, userMiddleware)
	get.HandleFunc("/users/{name:[[:alnum:]]+}", h.GetUser).Methods("GET", "OPTIONS")
	get.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", h.GetOutbox).Methods("GET", "OPTIONS")
	// get.HandleFunc("/users/{name:[[:alnum:]]+}/following", getFollowing).Methods("GET", "OPTIONS")
	// get.HandleFunc("/users/{name:[[:alnum:]]+}/followers", getFollowers).Methods("GET", "OPTIONS")
	// get.HandleFunc("/users/{name:[[:alnum:]]+}/liked", getLiked).Methods("GET", "OPTIONS")
	// get.HandleFunc("/activities/{id}", getActivity).Methods("GET", "OPTIONS")
	// get.HandleFunc("/objects/{id}", getObject).Methods("GET", "OPTIONS")

	// post := r.NewRoute().Subrouter() // -> public POST requests
	// post.Use(contentTypeMiddleware, userMiddleware)
	// post.HandleFunc("/users/{name:[[:alnum:]]+}/inbox", postInbox).Methods("POST", "OPTIONS")

	// aGet := get.NewRoute().Subrouter()
	// aGet.Use(jwtMiddleware)
	// aGet.HandleFunc("/users/{name:[[:alnum:]]+}/inbox", getInbox).Methods("GET", "OPTIONS")

	// aPost := post.NewRoute().Subrouter()
	// aPost.Use(jwtMiddleware)
	// aPost.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", postOutbox).Methods("POST", "OPTIONS")

	// sink := r.NewRoute().Subrouter() // -> sink to handle all other routes
	// sink.Use(acceptMiddleware)
	// sink.PathPrefix("/").HandlerFunc(sinkHandler).Methods("GET", "OPTIONS")
	return h
}

// TODO: Change this to return only what is needed for http.ListenAndServe(...)
func (h *MuxHandler) GetRouter() *mux.Router {
	return h.router
}

func (h *MuxHandler) AllowCORS(allowedOrigins []string) {
	cors := h.middleware.CreateCORSMiddleware(allowedOrigins)
	h.router.Use(cors)
}

func (h *MuxHandler) GetWebFinger(w http.ResponseWriter, r *http.Request) {
	resource := r.FormValue("resource")
	name, err := h.resource.ParseResource(resource)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	user, err := h.service.DiscoverUserByName(name)
	if err != nil {
		h.response.NotFound(w, err)
		return
	}
	webfinger := h.resource.GenerateWebFinger(user.Name)
	w.Header().Set("Content-Type", "application/jrd+json")
	json.NewEncoder(w).Encode(webfinger)
}

func (h *MuxHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := h.service.GetUserByName(name)
	if err != nil {
		h.response.NotFound(w, err)
		return
	}
	actor := h.resource.GenerateActor(user.Name)
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(actor)
}

func (h *MuxHandler) GetOutbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := h.service.GetUserByName(name)
	if err != nil {
		h.response.NotFound(w, err)
		return
	}
	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := h.service.GetOutboxTotalItemsByUserName(user.Name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		// TODO: abstract this out to the resource package?
		// would be a call to h.resource.GenerateOutbox
		// this would remove the config dependency (for now...)
		outbox := h.resource.GenerateOrderedCollection(user.Name, h.endpoints.Outbox, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(outbox)
		return
	}
	activities, err := h.service.GetOutboxByUserName(user.Name)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	// TODO: abstract this out to the resource package?
	// would be a call to h.resource.GenerateOutboxPage
	// this would remove the config dependency (for now...)
	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}
	outboxPage := h.resource.GenerateOrderedCollectionPage(name, h.endpoints.Outbox, orderedItems)
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(outboxPage)
}
