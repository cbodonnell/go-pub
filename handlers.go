package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

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
	w.Header().Set("Content-Type", "application/jrd+json")
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
	w.Header().Set("Content-Type", "application/jrd+json")
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
		inbox := generateOrderedCollection(user.Name, config.Endpoints.Inbox)
		w.Header().Set("Content-Type", "application/jrd+json")
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
	w.Header().Set("Content-Type", "application/jrd+json")
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
		outbox := generateOrderedCollection(user.Name, config.Endpoints.Outbox)
		w.Header().Set("Content-Type", "application/jrd+json")
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
	w.Header().Set("Content-Type", "application/jrd+json")
	json.NewEncoder(w).Encode(outboxPage)
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
		following := generateOrderedCollection(user.Name, config.Endpoints.Following)
		w.Header().Set("Content-Type", "application/jrd+json")
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
	w.Header().Set("Content-Type", "application/jrd+json")
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
		followers := generateOrderedCollection(user.Name, config.Endpoints.Followers)
		w.Header().Set("Content-Type", "application/jrd+json")
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
	w.Header().Set("Content-Type", "application/jrd+json")
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
		liked := generateOrderedCollection(user.Name, config.Endpoints.Liked)
		w.Header().Set("Content-Type", "application/jrd+json")
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
	w.Header().Set("Content-Type", "application/jrd+json")
	json.NewEncoder(w).Encode(likedPage)
}
