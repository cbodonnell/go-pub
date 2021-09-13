package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cheebz/go-pub/activitypub"
	"github.com/cheebz/go-pub/cache"
	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/middleware"
	"github.com/cheebz/go-pub/models"
	"github.com/cheebz/go-pub/resources"
	"github.com/cheebz/go-pub/responses"
	"github.com/cheebz/go-pub/services"
	"github.com/cheebz/go-pub/utils"
	"github.com/cheebz/sigs"
	"github.com/gorilla/mux"
)

type MuxHandler struct {
	endpoints  config.Endpoints
	middleware middleware.Middleware
	service    services.Service
	resource   resources.Resource
	response   responses.Response
	cache      cache.Cache
	router     *mux.Router
}

var (
	nameParam = "name"
)

func NewMuxHandler(_endpoints config.Endpoints, _middleware middleware.Middleware, _service services.Service, _resource resources.Resource, _response responses.Response, _cache cache.Cache) Handler {
	h := &MuxHandler{
		endpoints:  _endpoints,
		middleware: _middleware,
		service:    _service,
		resource:   _resource,
		response:   _response,
		cache:      _cache,
		router:     mux.NewRouter(),
	}
	h.setupRoutes()
	return h
}

func (h *MuxHandler) setupRoutes() {
	wf := h.router.NewRoute().Subrouter() // -> webfinger
	wf.HandleFunc("/.well-known/webfinger", h.GetWebFinger).Methods("GET", "OPTIONS")

	userMiddleware := h.middleware.CreateUserMiddleware(h.service)

	get := h.router.NewRoute().Subrouter() // -> public GET requests
	get.Use(h.middleware.AcceptMiddleware, userMiddleware)
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}", h.endpoints.Users, nameParam), h.GetUser).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.endpoints.Users, nameParam, h.endpoints.Outbox), h.GetOutbox).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.endpoints.Users, nameParam, h.endpoints.Following), h.GetFollowing).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.endpoints.Users, nameParam, h.endpoints.Followers), h.GetFollowers).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.endpoints.Users, nameParam, h.endpoints.Liked), h.GetLiked).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{id}", h.endpoints.Activities), h.GetActivity).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{id}", h.endpoints.Objects), h.GetObject).Methods("GET", "OPTIONS")

	post := h.router.NewRoute().Subrouter() // -> public POST requests
	post.Use(h.middleware.ContentTypeMiddleware, userMiddleware)
	post.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.endpoints.Users, nameParam, h.endpoints.Inbox), h.PostInbox).Methods("POST", "OPTIONS")

	jwtUsernameMiddleware := h.middleware.CreateJwtUsernameMiddleware(nameParam)

	aGet := get.NewRoute().Subrouter()
	aGet.Use(jwtUsernameMiddleware)
	aGet.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.endpoints.Users, nameParam, h.endpoints.Inbox), h.GetInbox).Methods("GET", "OPTIONS")

	aPost := post.NewRoute().Subrouter()
	aPost.Use(jwtUsernameMiddleware)
	aPost.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.endpoints.Users, nameParam, h.endpoints.Outbox), h.PostOutbox).Methods("POST", "OPTIONS")

	sink := h.router.NewRoute().Subrouter() // -> sink to handle all other routes
	sink.Use(h.middleware.AcceptMiddleware)
	sink.PathPrefix("/").HandlerFunc(h.SinkHandler).Methods("GET", "OPTIONS")

	// Static files
	// TODO: Static files
	// r.PathPrefix(fmt.Sprintf("/files/{%s:[[:alnum:]]+}/", nameParam)).Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("./static/"))))
	// r.Use(jwtUsernameMiddleware)
}

// TODO: Change this to return only what is needed for http.ListenAndServe(...)
func (h *MuxHandler) GetRouter() http.Handler {
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
	webfinger, err := h.cache.Get(fmt.Sprintf("webfinger-%s", name), models.WebFinger{})
	if err != nil {
		log.Println(fmt.Sprintf("no cache for webfinger-%s", name))
		user, err := h.service.DiscoverUserByName(name)
		if err != nil {
			h.response.NotFound(w, err)
			return
		}
		webfinger = h.resource.GenerateWebFinger(user.Name)
		err = h.cache.Set(fmt.Sprintf("webfinger-%s", name), webfinger)
		if err != nil {
			log.Println(fmt.Sprintf("failed to cache webfinger-%s: %+v", name, webfinger))
		}
	}
	w.Header().Set("Content-Type", "application/jrd+json")
	json.NewEncoder(w).Encode(webfinger)
}

func (h *MuxHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	actor, err := h.cache.Get(fmt.Sprintf("actor-%s", name), models.Actor{})
	if err != nil {
		log.Println(fmt.Sprintf("no cache for actor-%s", name))
		user, err := h.service.GetUserByName(name)
		if err != nil {
			h.response.NotFound(w, err)
			return
		}
		actor = h.resource.GenerateActor(user.Name)
		err = h.cache.Set(fmt.Sprintf("actor-%s", name), actor)
		if err != nil {
			log.Println(fmt.Sprintf("failed to cache actor-%s: %+v", name, actor))
		}
	}
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(actor)
}

func (h *MuxHandler) GetInbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := h.service.GetInboxTotalItemsByUserName(name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		inbox := h.resource.GenerateOrderedCollection(name, h.endpoints.Inbox, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(inbox)
		return
	}
	activities, err := h.service.GetInboxByUserName(name)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}
	inboxPage := h.resource.GenerateOrderedCollectionPage(name, h.endpoints.Inbox, orderedItems)
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(inboxPage)
}

func (h *MuxHandler) GetOutbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
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
	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}
	outboxPage := h.resource.GenerateOrderedCollectionPage(name, h.endpoints.Outbox, orderedItems)
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(outboxPage)
}

func (h *MuxHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	user, err := h.service.GetUserByName(name)
	if err != nil {
		h.response.NotFound(w, err)
		return
	}
	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := h.service.GetFollowingTotalItemsByUserName(user.Name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		following := h.resource.GenerateOrderedCollection(user.Name, h.endpoints.Following, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(following)
		return
	}
	following, err := h.service.GetFollowingByUserName(user.Name)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	orderedItems := make([]interface{}, len(following))
	for i, actor := range following {
		orderedItems[i] = actor
	}
	followingPage := h.resource.GenerateOrderedCollectionPage(user.Name, h.endpoints.Following, orderedItems)
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(followingPage)
}

func (h *MuxHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	user, err := h.service.GetUserByName(name)
	if err != nil {
		h.response.NotFound(w, err)
		return
	}
	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := h.service.GetFollowersTotalItemsByUserName(user.Name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		followers := h.resource.GenerateOrderedCollection(user.Name, h.endpoints.Followers, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(followers)
		return
	}
	followers, err := h.service.GetFollowersByUserName(user.Name)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	orderedItems := make([]interface{}, len(followers))
	for i, actor := range followers {
		orderedItems[i] = actor
	}
	followersPage := h.resource.GenerateOrderedCollectionPage(user.Name, h.endpoints.Followers, orderedItems)
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(followersPage)
}

func (h *MuxHandler) GetLiked(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	user, err := h.service.GetUserByName(name)
	if err != nil {
		h.response.NotFound(w, err)
		return
	}
	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := h.service.GetLikedTotalItemsByUserName(user.Name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		liked := h.resource.GenerateOrderedCollection(user.Name, h.endpoints.Liked, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(liked)
		return
	}
	liked, err := h.service.GetLikedByUserName(user.Name)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	orderedItems := make([]interface{}, len(liked))
	for i, activity := range liked {
		orderedItems[i] = activity
	}
	likedPage := h.resource.GenerateOrderedCollectionPage(user.Name, h.endpoints.Liked, orderedItems)
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(likedPage)
}

func (h *MuxHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	activity, err := h.service.GetActivity(id)
	if err != nil {
		h.response.NotFound(w, err)
		return
	}
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(activity)
}

func (h *MuxHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	object, err := h.service.GetObject(id)
	if err != nil {
		h.response.NotFound(w, err)
		return
	}
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(object)
}

func (h *MuxHandler) PostInbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	err := h.service.CheckUser(name)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	err = activitypub.CheckContentType(r.Header)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	payload, err := utils.ParseLimitedPayload(r.Body, 1*1024*1024)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	_, err = sigs.VerifyRequest(r, payload, activitypub.FetchPublicKeyString)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	activityArb, err := activitypub.ParsePayload(payload)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	activityArb, err = h.service.SaveInboxActivity(activityArb, name)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	h.response.Accepted(w)
}

func (h *MuxHandler) PostOutbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	// TODO: Roll this into a CheckJWTUsername middleware
	// claims, _ := jwt.CheckJWTClaims(r)
	// if claims.Username != name {
	// 	h.response.UnauthorizedRequest(w, errors.New("not your outbox"))
	// 	return
	// }
	err := activitypub.CheckContentType(r.Header)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	payload, err := utils.ParseLimitedPayload(r.Body, 1*1024*1024)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	activityArb, err := activitypub.ParsePayload(payload)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	activityArb, err = h.service.SaveOutboxActivity(activityArb, name)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", activitypub.ContentType)
	iri, err := activityArb.GetString("id")
	h.response.Created(w, iri)
	activityArb.Write(w)
}

func (h *MuxHandler) SinkHandler(w http.ResponseWriter, r *http.Request) {
	h.response.NotFound(w, errors.New("endpoint does not exist"))
	return
}
