package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cheebz/arb"
	"github.com/gorilla/mux"
)

var accept = "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\""
var acceptHeaders = http.Header{
	"Accept": []string{
		"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
		"application/activity+json",
	},
}

var contentType = "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\""
var contentTypeHeaders = http.Header{
	"Content-Type": []string{
		"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
		"application/activity+json",
	},
}

func checkContentType(headers http.Header) error {
	h := headers.Values("Content-Type")
	for _, v := range h {
		fmt.Println("Request contains Content-Type header: " + v)
		for _, item := range contentTypeHeaders["Content-Type"] {
			if strings.Contains(v, item) {
				return nil
			}
		}
	}
	return errors.New("invalid content-type headers")
}

func checkAccept(headers http.Header) error {
	h := headers.Values("Accept")
	for _, v := range h {
		fmt.Println("Request contains Accept header: " + v)
		for _, item := range acceptHeaders["Accept"] {
			if strings.Contains(v, item) {
				return nil
			}
		}
	}
	return errors.New("invalid accept headers")
}

func sinkHandler(w http.ResponseWriter, r *http.Request) {
	notFound(w, errors.New("endpoint does not exist"))
	return
}

func getWebFinger(w http.ResponseWriter, r *http.Request) {
	resource := r.FormValue("resource")
	name, err := parseResource(resource)
	if err != nil {
		badRequest(w, err)
		return
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
	w.Header().Set("Content-Type", contentType)
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
		// totalItems, err := queryInboxTotalItemsByUserName(user.Name)
		// if err != nil {
		// 	internalServerError(w, err)
		// 	return
		// }
		totalItems := 1

		inbox := generateOrderedCollection(user.Name, config.Endpoints.Inbox, totalItems)
		w.Header().Set("Content-Type", contentType)
		json.NewEncoder(w).Encode(inbox)
		return
	}

	// posts, err := queryInboxByUserName(user.Name)
	// if err != nil {
	// 	internalServerError(w, err)
	// 	return
	// }

	activities := make([]Activity, 0)

	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}

	inboxPage := generateOrderedCollectionPage(name, config.Endpoints.Inbox, orderedItems)
	w.Header().Set("Content-Type", contentType)
	json.NewEncoder(w).Encode(inboxPage)
}

func postInbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	err := checkUser(name)
	if err != nil {
		badRequest(w, err)
		return
	}
	recipient := fmt.Sprintf("%s://%s/%s/%s", config.Protocol, config.ServerName, config.Endpoints.Users, name)
	err = checkContentType(r.Header)
	if err != nil {
		badRequest(w, err)
		return
	}
	// TODO: Check signature here...
	payloadArb, err := arb.Read(r.Body)
	if err != nil {
		badRequest(w, err)
		return
	}
	err = checkContext(payloadArb)
	if err != nil {
		badRequest(w, err)
		return
	}
	payloadType, err := payloadArb.GetString("type")
	if err != nil {
		badRequest(w, err)
		return
	}
	var activityArb arb.Arb
	if !isActivity(payloadType) {
		badRequest(w, err)
		return
	}
	activityArb = payloadArb
	actorArb, err := findProp(activityArb, "actor", acceptHeaders)
	if err != nil {
		badRequest(w, err)
		return
	}
	actorIRI, err := actorArb.GetString("id")
	if err != nil {
		badRequest(w, err)
		return
	}
	objectArb, err := findProp(activityArb, "object", acceptHeaders)
	if err != nil {
		badRequest(w, err)
		return
	}
	objectIRI, err := objectArb.GetString("id")
	if err != nil {
		badRequest(w, err)
		return
	}
	if objectIRI != recipient {
		badRequest(w, errors.New("wrong inbox"))
		return
	}
	activityType, err := activityArb.GetString("type")
	if err != nil {
		badRequest(w, err)
		return
	}
	switch activityType {
	case "Follow":
		activityArb, err = createInboxActivity(activityArb, recipient, actorIRI, recipient)
		if err != nil {
			internalServerError(w, err)
			return
		}
	default:
		badRequest(w, errors.New("unsupported activity type"))
		return
		// Activity type is something else, save object reference (if new), Activity, and Activity_to
	}

	for k, l := range contentTypeHeaders {
		for _, v := range l {
			w.Header().Add(k, v)
		}
	}
	activityArb.Write(w)
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
		w.Header().Set("Content-Type", contentType)
		json.NewEncoder(w).Encode(outbox)
		return
	}

	activities, err := queryOutboxByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}

	outboxPage := generateOrderedCollectionPage(name, config.Endpoints.Outbox, orderedItems)
	w.Header().Set("Content-Type", contentType)
	json.NewEncoder(w).Encode(outboxPage)
}

func postOutbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	claims, _ := checkJWTClaims(r)
	if claims.Username != name {
		unauthorizedRequest(w, errors.New("not your outbox"))
		return
	}
	err := checkContentType(r.Header)
	if err != nil {
		badRequest(w, err)
		return
	}
	payloadArb, err := arb.Read(r.Body)
	if err != nil {
		badRequest(w, err)
		return
	}
	err = checkContext(payloadArb)
	if err != nil {
		badRequest(w, err)
		return
	}
	payloadType, err := payloadArb.GetString("type")
	if err != nil {
		badRequest(w, err)
		return
	}
	// TODO: Refactor into a parsePayload method
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
	activityType, err := activityArb.GetString("type")
	if err != nil {
		badRequest(w, err)
		return
	}
	actor := fmt.Sprintf("%s://%s/%s/%s", config.Protocol, config.ServerName, config.Endpoints.Users, claims.Username)
	switch activityType {
	case "Create":
		objectArb, err := activityArb.GetArb("object")
		if err != nil {
			badRequest(w, err)
			return
		}
		activityArb, err = createOutboxActivity(activityArb, objectArb, actor)
		if err != nil {
			internalServerError(w, err)
			return
		}
	case "Like":
		activityArb, err = createOutboxReferenceActivity(activityArb, actor)
		if err != nil {
			internalServerError(w, err)
			return
		}
	default:
		badRequest(w, errors.New("unsupported activity type"))
		return
		// Activity type is something else, save object reference (if new), Activity, and Activity_to
	}

	// TODO: Propagate Activity <-- Can this be done with a concurrent worker
	// by passing the activity into a channel?

	for k, l := range contentTypeHeaders {
		for _, v := range l {
			w.Header().Add(k, v)
		}
	}
	iri, err := activityArb.GetString("id")
	created(w, iri)
	activityArb.Write(w)
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
		w.Header().Set("Content-Type", contentType)
		json.NewEncoder(w).Encode(following)
		return
	}

	// // TODO: Implement a method to get the following collection
	// activities, err := queryOutboxByUserName(user.Name)
	// if err != nil {
	// 	internalServerError(w, err)
	// 	return
	// }

	activities := make([]Activity, 0)

	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}

	followingPage := generateOrderedCollectionPage(user.Name, config.Endpoints.Following, orderedItems)
	w.Header().Set("Content-Type", contentType)
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
		w.Header().Set("Content-Type", contentType)
		json.NewEncoder(w).Encode(followers)
		return
	}

	// // TODO: Implement a method to get the followers collection
	// activities, err := queryOutboxByUserName(user.Name)
	// if err != nil {
	// 	internalServerError(w, err)
	// 	return
	// }

	activities := make([]Activity, 0)

	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}

	followersPage := generateOrderedCollectionPage(user.Name, config.Endpoints.Followers, orderedItems)
	w.Header().Set("Content-Type", contentType)
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
		w.Header().Set("Content-Type", contentType)
		json.NewEncoder(w).Encode(liked)
		return
	}

	// // TODO: Implement a method to get the liked collection
	// activities, err := queryOutboxByUserName(user.Name)
	// if err != nil {
	// 	internalServerError(w, err)
	// 	return
	// }

	activities := make([]Activity, 0)

	orderedItems := make([]interface{}, len(activities))
	for i, activity := range activities {
		orderedItems[i] = activity
	}

	likedPage := generateOrderedCollectionPage(user.Name, config.Endpoints.Liked, orderedItems)
	w.Header().Set("Content-Type", contentType)
	json.NewEncoder(w).Encode(likedPage)
}

func getActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		badRequest(w, err)
		return
	}

	activity, err := queryActivity(id)
	if err != nil {
		notFound(w, err)
		return
	}

	w.Header().Set("Content-Type", contentType)
	json.NewEncoder(w).Encode(activity)
}

func getObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		badRequest(w, err)
		return
	}

	object, err := queryObject(id)
	if err != nil {
		notFound(w, err)
		return
	}

	w.Header().Set("Content-Type", contentType)
	json.NewEncoder(w).Encode(object)
}
