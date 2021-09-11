package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/middleware"
	"github.com/cheebz/go-pub/models"
	"github.com/cheebz/go-pub/repositories"
	"github.com/cheebz/go-pub/responses"
	"github.com/gorilla/mux"
)

type MuxHandler struct {
	repo repositories.Repository
}

var (
	conf          config.Configuration
	repo          repositories.Repository
	accept        = "application/activity+json"
	acceptHeaders = http.Header{
		"Accept": []string{
			"application/activity+json",
			"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
		},
	}
	contentType        = "application/activity+json"
	contentTypeHeaders = http.Header{
		"Content-Type": []string{
			"application/activity+json",
			"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
		},
	}
)

func NewMuxHandler(_conf config.Configuration, _repo repositories.Repository) *mux.Router {
	h := &MuxHandler{
		repo: repo,
	}
	conf = _conf
	repo = _repo
	r := mux.NewRouter()

	// wf := r.NewRoute().Subrouter() // -> webfinger
	// wf.HandleFunc("/.well-known/webfinger", getWebFinger).Methods("GET", "OPTIONS")

	userMiddleware := middleware.CreateUserMiddleware(repo)

	get := r.NewRoute().Subrouter() // -> public GET requests
	get.Use(middleware.AcceptMiddleware, userMiddleware)
	get.HandleFunc("/users/{name:[[:alnum:]]+}", h.GetUser).Methods("GET", "OPTIONS")
	// get.HandleFunc("/users/{name:[[:alnum:]]+}/outbox", getOutbox).Methods("GET", "OPTIONS")
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
	return r
}

func (h *MuxHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := repo.QueryUserByName(name)
	if err != nil {
		responses.NotFound(w, err)
		return
	}

	actor := generateActor(user.Name)
	w.Header().Set("Content-Type", contentType)
	json.NewEncoder(w).Encode(actor)
}

func generateActor(name string) models.Actor {
	return models.Actor{
		Object: models.Object{
			Context: []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
				map[string]interface{}{
					"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
				},
			},
			Id:      fmt.Sprintf("%s://%s/%s/%s", conf.Protocol, conf.ServerName, conf.Endpoints.Users, name),
			Type:    "Person",
			Name:    name,
			Url:     fmt.Sprintf("%s://%s/%s/%s", conf.Protocol, conf.ServerName, conf.Endpoints.Users, name),
			Summary: fmt.Sprintf("Summary of %s to come...", name), // TODO: Implement this
		},
		Inbox:                     fmt.Sprintf("%s://%s/%s/%s/%s", conf.Protocol, conf.ServerName, conf.Endpoints.Users, name, conf.Endpoints.Inbox),
		Outbox:                    fmt.Sprintf("%s://%s/%s/%s/%s", conf.Protocol, conf.ServerName, conf.Endpoints.Users, name, conf.Endpoints.Outbox),
		Following:                 fmt.Sprintf("%s://%s/%s/%s/%s", conf.Protocol, conf.ServerName, conf.Endpoints.Users, name, conf.Endpoints.Following),
		Followers:                 fmt.Sprintf("%s://%s/%s/%s/%s", conf.Protocol, conf.ServerName, conf.Endpoints.Users, name, conf.Endpoints.Followers),
		Liked:                     fmt.Sprintf("%s://%s/%s/%s/%s", conf.Protocol, conf.ServerName, conf.Endpoints.Users, name, conf.Endpoints.Liked),
		PreferredUsername:         name,
		ManuallyApprovesFollowers: false, // TODO: Implement this
		PublicKey: models.PublicKey{
			ID:           fmt.Sprintf("%s://%s/%s/%s#main-key", conf.Protocol, conf.ServerName, conf.Endpoints.Users, name),
			Owner:        fmt.Sprintf("%s://%s/%s/%s", conf.Protocol, conf.ServerName, conf.Endpoints.Users, name),
			PublicKeyPem: conf.RSAPublicKey,
		},
	}
}
