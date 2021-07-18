package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/cheebz/arb"
	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("static/templates/*.html"))
var acceptHeaders = http.Header{
	"Accept": []string{
		"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
		"application/activity+json",
	},
}
var contentTypeHeaders = http.Header{
	"Content-Type": []string{
		"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
		"application/activity+json",
	},
}

func renderTemplate(w http.ResponseWriter, template string, data interface{}) {
	err := templates.ExecuteTemplate(w, template, data)
	if err != nil {
		internalServerError(w, err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	claims, _ := checkJWTClaims(r)
	data := HomeData{
		Claims:         claims,
		ServerName:     config.ServerName,
		UsersEndpoint:  config.Endpoints.Users,
		OutboxEndpoint: config.Endpoints.Outbox,
		Auth:           config.Auth,
	}
	renderTemplate(w, "index.html", data)
}

func getWebFinger(w http.ResponseWriter, r *http.Request) {
	resource := r.FormValue("resource")
	name, err := parseResource(resource)
	if err != nil {
		badRequest(w, err)
	}

	user, err := queryUserByName(name)
	if err != nil {
		notFound(w, err)
		return
	}

	if !user.Discoverable {
		notFound(w, errors.New("user is not discoverable"))
		return
	}

	webfinger := generateWebFinger(user.Name)
	w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	json.NewEncoder(w).Encode(webfinger)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := queryUserByName(name)
	if err != nil {
		notFound(w, err)
		return
	}

	actor := generateActor(user.Name)
	w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	json.NewEncoder(w).Encode(actor)
}

func getInbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := queryUserByName(name)
	if err != nil {
		notFound(w, err)
		return
	}

	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := queryInboxTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		inbox := generateOrderedCollection(user.Name, config.Endpoints.Inbox, totalItems)
		w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
		json.NewEncoder(w).Encode(inbox)
		return
	}

	posts, err := queryInboxByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

	orderedItems := make([]interface{}, len(posts))
	for i, post := range posts {
		orderedItems[i] = generatePostActivity(post)
	}

	inboxPage := generateOrderedCollectionPage(name, config.Endpoints.Inbox, orderedItems)
	w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	json.NewEncoder(w).Encode(inboxPage)
}

func getOutbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := queryUserByName(name)
	if err != nil {
		notFound(w, err)
		return
	}

	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := queryOutboxTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		outbox := generateOrderedCollection(user.Name, config.Endpoints.Outbox, totalItems)
		w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
		json.NewEncoder(w).Encode(outbox)
		return
	}

	posts, err := queryOutboxByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

	orderedItems := make([]interface{}, len(posts))
	for i, post := range posts {
		orderedItems[i] = generatePostActivity(post)
	}

	outboxPage := generateOrderedCollectionPage(name, config.Endpoints.Outbox, orderedItems)
	w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	json.NewEncoder(w).Encode(outboxPage)
}

func postOutbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	claims, _ := checkJWTClaims(r)
	if claims.Username != name {
		unauthorizedRequest(w, errors.New("not your outbox"))
		return
	}
	payloadArb, err := arb.Read(r.Body)
	if err != nil {
		badRequest(w, err)
		return
	}
	payloadType, err := payloadArb.GetString("type")
	if err != nil {
		badRequest(w, err)
		return
	}
	// fmt.Println(fmt.Sprintf("Payload of type %s", payloadType))
	var activityArb arb.Arb
	if isObject(payloadType) {
		activityArb, err = createActivity(payloadArb)
		if err != nil {
			badRequest(w, err)
			return
		}
	}
	if isActivity(payloadType) {
		activityArb = payloadArb
		objectArb, err := findObject(activityArb, acceptHeaders)
		err = objectArb.PropToArray("@context")
		if err != nil {
			badRequest(w, err)
			return
		}
		err = formatRecipients(objectArb)
		if err != nil {
			badRequest(w, err)
			return
		}
		activityArb["object"] = objectArb
		if err != nil {
			badRequest(w, err)
			return
		}
		err = activityArb.PropToArray("@context")
		if err != nil {
			badRequest(w, err)
			return
		}
		err = formatRecipients(activityArb)
		if err != nil {
			badRequest(w, err)
			return
		}
	}
	if activityArb == nil {
		badRequest(w, err)
		return
	}
	var activity ActivityResource
	err = json.Unmarshal(activityArb.ToBytes(), &activity)
	if err != nil {
		badRequest(w, err)
		return
	}
	switch payloadType {
	case "Create":
		// Activity type is Create, save object detail, Activity_to, and Activity
		// set attributedTo?
		// Set object IRI

	default:
		// Activity type is something else, save object reference (if new), Activity_to, and Activity
	}

	// Apply generated ID
	// Apply actor
	// Propagate Activity <-- Can this be done with a worker?
	// Resolve addressing between object and activity using to, bto, cc, bcc, and audience

	for k, l := range contentTypeHeaders {
		for _, v := range l {
			w.Header().Add(k, v)
		}
	}
	created(w, activity.Id)
	json.NewEncoder(w).Encode(activity)
}

func getFollowing(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := queryUserByName(name)
	if err != nil {
		notFound(w, err)
		return
	}

	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := queryOutboxTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		following := generateOrderedCollection(user.Name, config.Endpoints.Following, totalItems)
		w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
		json.NewEncoder(w).Encode(following)
		return
	}

	// TODO: Implement a method to get the following collection
	posts, err := queryOutboxByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

	orderedItems := make([]interface{}, len(posts))
	for i, post := range posts {
		orderedItems[i] = generatePostActivity(post)
	}

	followingPage := generateOrderedCollectionPage(user.Name, config.Endpoints.Following, orderedItems)
	w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	json.NewEncoder(w).Encode(followingPage)
}

func getFollowers(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := queryUserByName(name)
	if err != nil {
		notFound(w, err)
		return
	}

	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := queryOutboxTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		followers := generateOrderedCollection(user.Name, config.Endpoints.Followers, totalItems)
		w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
		json.NewEncoder(w).Encode(followers)
		return
	}

	// TODO: Implement a method to get the followers collection
	posts, err := queryOutboxByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

	orderedItems := make([]interface{}, len(posts))
	for i, post := range posts {
		orderedItems[i] = generatePostActivity(post)
	}

	followersPage := generateOrderedCollectionPage(user.Name, config.Endpoints.Followers, orderedItems)
	w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	json.NewEncoder(w).Encode(followersPage)
}

func getLiked(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := queryUserByName(name)
	if err != nil {
		notFound(w, err)
		return
	}

	page := r.FormValue("page")
	if page != "true" {
		totalItems, err := queryOutboxTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		liked := generateOrderedCollection(user.Name, config.Endpoints.Liked, totalItems)
		w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
		json.NewEncoder(w).Encode(liked)
		return
	}

	// TODO: Implement a method to get the liked collection
	posts, err := queryOutboxByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

	orderedItems := make([]interface{}, len(posts))
	for i, post := range posts {
		orderedItems[i] = generatePostActivity(post)
	}

	likedPage := generateOrderedCollectionPage(user.Name, config.Endpoints.Liked, orderedItems)
	w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
	json.NewEncoder(w).Encode(likedPage)
}
