package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cheebz/go-pub/activitypub"
	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/media"
	"github.com/cheebz/go-pub/middleware"
	"github.com/cheebz/go-pub/resources"
	"github.com/cheebz/go-pub/responses"
	"github.com/cheebz/go-pub/services"
	"github.com/cheebz/go-pub/utils"
	"github.com/cheebz/sigs"
	"github.com/gorilla/mux"
)

type MuxHandler struct {
	conf       config.Configuration
	middleware middleware.Middleware
	service    services.Service
	resource   resources.Resource
	response   responses.Response
	router     *mux.Router
}

var (
	nameParam = "name"
)

func NewMuxHandler(_config config.Configuration, _middleware middleware.Middleware, _service services.Service, _resource resources.Resource, _response responses.Response) Handler {
	h := &MuxHandler{
		conf:       _config,
		middleware: _middleware,
		service:    _service,
		resource:   _resource,
		response:   _response,
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
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}", h.conf.Endpoints.Users, nameParam), h.GetUser).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.conf.Endpoints.Users, nameParam, h.conf.Endpoints.Outbox), h.GetOutbox).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.conf.Endpoints.Users, nameParam, h.conf.Endpoints.Following), h.GetFollowing).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.conf.Endpoints.Users, nameParam, h.conf.Endpoints.Followers), h.GetFollowers).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.conf.Endpoints.Users, nameParam, h.conf.Endpoints.Liked), h.GetLiked).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{id}", h.conf.Endpoints.Activities), h.GetActivity).Methods("GET", "OPTIONS")
	get.HandleFunc(fmt.Sprintf("/%s/{id}", h.conf.Endpoints.Objects), h.GetObject).Methods("GET", "OPTIONS")

	post := h.router.NewRoute().Subrouter() // -> public POST requests
	post.Use(h.middleware.ContentTypeMiddleware, userMiddleware)
	post.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.conf.Endpoints.Users, nameParam, h.conf.Endpoints.Inbox), h.PostInbox).Methods("POST", "OPTIONS")

	jwtUsernameMiddleware := h.middleware.CreateJwtUsernameMiddleware(nameParam)

	aGet := get.NewRoute().Subrouter()
	aGet.Use(jwtUsernameMiddleware)
	aGet.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.conf.Endpoints.Users, nameParam, h.conf.Endpoints.Inbox), h.GetInbox).Methods("GET", "OPTIONS")

	aPost := post.NewRoute().Subrouter()
	aPost.Use(jwtUsernameMiddleware)
	aPost.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.conf.Endpoints.Users, nameParam, h.conf.Endpoints.Outbox), h.PostOutbox).Methods("POST", "OPTIONS")

	uPost := h.router.NewRoute().Subrouter() // -> authenticated uploads POST
	uPost.Use(jwtUsernameMiddleware)
	uPost.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.conf.Endpoints.Users, nameParam, h.conf.Endpoints.UploadMedia), h.UploadMedia).Methods("POST", "OPTIONS")

	uGet := h.router.NewRoute().Subrouter() // -> authenticated uploads GET
	uGet.PathPrefix(fmt.Sprintf("/%s/", h.conf.Endpoints.Uploads)).Handler(http.StripPrefix(fmt.Sprintf("/%s/", h.conf.Endpoints.Uploads), http.FileServer(http.Dir(h.conf.UploadDir))))

	cGet := h.router.NewRoute().Subrouter() // -> authenticated checks GET
	uPost.Use(jwtUsernameMiddleware)
	cGet.HandleFunc(fmt.Sprintf("/%s/{%s:[[:alnum:]]+}/%s", h.conf.Endpoints.Users, nameParam, h.conf.Endpoints.Check), h.CheckActivity).Methods("GET", "OPTIONS")

	sink := h.router.NewRoute().Subrouter() // -> sink to handle all other routes
	sink.Use(h.middleware.AcceptMiddleware)
	sink.PathPrefix("/").HandlerFunc(h.SinkHandler).Methods("GET", "OPTIONS")

}

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
	name := mux.Vars(r)[nameParam]
	user, err := h.service.GetUserByName(name)
	if err != nil {
		h.response.NotFound(w, err)
		return
	}
	actor := h.resource.GenerateActor(user.Name)
	w.Header().Set("Content-Type", activitypub.ContentType)
	json.NewEncoder(w).Encode(actor)
}

func (h *MuxHandler) GetInbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	page := r.FormValue("page")
	if page == "" {
		totalItems, err := h.service.GetInboxTotalItemsByUserName(name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		inbox := h.resource.GenerateOrderedCollection(name, h.conf.Endpoints.Inbox, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(inbox)
		return
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	activities, err := h.service.GetInboxByUserName(name, pageNum)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}
	inboxPage := h.resource.GenerateOrderedCollectionPage(name, h.conf.Endpoints.Inbox, orderedItems, pageNum)
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
	if page == "" {
		totalItems, err := h.service.GetOutboxTotalItemsByUserName(user.Name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		outbox := h.resource.GenerateOrderedCollection(user.Name, h.conf.Endpoints.Outbox, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(outbox)
		return
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	activities, err := h.service.GetOutboxByUserName(user.Name, pageNum)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}
	outboxPage := h.resource.GenerateOrderedCollectionPage(name, h.conf.Endpoints.Outbox, orderedItems, pageNum)
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
	if page == "" {
		totalItems, err := h.service.GetFollowingTotalItemsByUserName(user.Name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		following := h.resource.GenerateOrderedCollection(user.Name, h.conf.Endpoints.Following, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(following)
		return
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	following, err := h.service.GetFollowingByUserName(user.Name, pageNum)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	orderedItems := make([]interface{}, len(following))
	for i, actor := range following {
		orderedItems[i] = actor
	}
	followingPage := h.resource.GenerateOrderedCollectionPage(user.Name, h.conf.Endpoints.Following, orderedItems, pageNum)
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
	if page == "" {
		totalItems, err := h.service.GetFollowersTotalItemsByUserName(user.Name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		followers := h.resource.GenerateOrderedCollection(user.Name, h.conf.Endpoints.Followers, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(followers)
		return
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	followers, err := h.service.GetFollowersByUserName(user.Name, pageNum)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	orderedItems := make([]interface{}, len(followers))
	for i, actor := range followers {
		orderedItems[i] = actor
	}
	followersPage := h.resource.GenerateOrderedCollectionPage(user.Name, h.conf.Endpoints.Followers, orderedItems, pageNum)
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
	if page == "" {
		totalItems, err := h.service.GetLikedTotalItemsByUserName(user.Name)
		if err != nil {
			h.response.InternalServerError(w, err)
			return
		}
		liked := h.resource.GenerateOrderedCollection(user.Name, h.conf.Endpoints.Liked, totalItems)
		w.Header().Set("Content-Type", activitypub.ContentType)
		json.NewEncoder(w).Encode(liked)
		return
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	liked, err := h.service.GetLikedByUserName(user.Name, pageNum)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	orderedItems := make([]interface{}, len(liked))
	for i, activity := range liked {
		orderedItems[i] = activity
	}
	likedPage := h.resource.GenerateOrderedCollectionPage(user.Name, h.conf.Endpoints.Liked, orderedItems, pageNum)
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
	_, err = h.service.SaveInboxActivity(activityArb, name)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	h.response.Accepted(w)
}

func (h *MuxHandler) PostOutbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
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
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	h.response.Created(w, iri)
	activityArb.Write(w)
}

func (h *MuxHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	err := activitypub.CheckUploadContentType(r.Header)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.response.BadRequest(w, err)
		return
	}
	activityArb, err := activitypub.ParsePayload([]byte(r.FormValue("object")))
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	file, err := media.ParseMedia(r, "file")
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}
	activityArb, err = h.service.UploadMedia(activityArb, file, name)
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", activitypub.ContentType)
	iri, err := activityArb.GetString("id")
	if err != nil {
		h.response.InternalServerError(w, err)
		return
	}
	h.response.Created(w, iri)
	activityArb.Write(w)
}

func (h *MuxHandler) CheckActivity(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)[nameParam]
	activityType := r.FormValue("activity")
	objectIRI := r.FormValue("object")
	activityIRI := h.service.CheckActivity(name, activityType, objectIRI)
	checkResponse := h.resource.GenerateCheckResponse(activityIRI)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkResponse)
}

func (h *MuxHandler) SinkHandler(w http.ResponseWriter, r *http.Request) {
	h.response.NotFound(w, fmt.Errorf("endpoint %s does not exist", r.URL))
}
