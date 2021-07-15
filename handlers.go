package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("static/templates/*.html"))

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
		totalItems, err := queryInboxTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		inbox := generateOrderedCollection(user.Name, config.Endpoints.Inbox, totalItems)
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
		totalItems, err := queryOutboxTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		outbox := generateOrderedCollection(user.Name, config.Endpoints.Outbox, totalItems)
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

func postOutbox(w http.ResponseWriter, r *http.Request) {

	// JUST TESTING STUFF...
	birdJson := `{"birds":[{"pigeon":"likes to perch on rocks"},{"eagle":"bird of prey"}],"animals":"none"}`
	var result map[string]interface{}
	json.Unmarshal([]byte(birdJson), &result)

	// TODO: Create a method to check types
	birds := result["birds"]
	switch c := birds.(type) {
	case string:
		fmt.Println("birds is string")
	case map[string]interface{}:
		fmt.Println("birds is map[string]interface{}")
	case []interface{}:
		fmt.Println("birds is []interface{}")
	default:
		fmt.Println(fmt.Sprintf("birds is might be %T", c))
	}
	animals := result["animals"]
	switch c := animals.(type) {
	case string:
		fmt.Println("animals is string")
	case map[string]interface{}:
		fmt.Println("animals is map[string]interface{}")
	case []interface{}:
		fmt.Println("animals is []interface{}")
	default:
		fmt.Println(fmt.Sprintf("animals is might be %T", c))
	}

	iri := "https://localhost/posts/123123"
	created(w, iri)
	return
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
		totalItems, err := queryOutboxTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		followers := generateOrderedCollection(user.Name, config.Endpoints.Followers, totalItems)
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
		totalItems, err := queryOutboxTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		liked := generateOrderedCollection(user.Name, config.Endpoints.Liked, totalItems)
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
