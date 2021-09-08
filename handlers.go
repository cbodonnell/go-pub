package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/cheebz/go-pub/repositories"
	"github.com/cheebz/sigs"
	"github.com/gorilla/mux"
)

var (
	repository repositories.Repository
)

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
	user, err := repository.QueryUserByName(name)
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
	claims, _ := checkJWTClaims(r)
	if claims.Username != name {
		unauthorizedRequest(w, errors.New("not your inbox"))
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
		w.Header().Set("Content-Type", contentType)
		json.NewEncoder(w).Encode(inbox)
		return
	}

	activities, err := queryInboxByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

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
	var buf bytes.Buffer
	limiter := io.LimitReader(r.Body, 1*1024*1024)
	io.Copy(&buf, limiter)
	payload := buf.Bytes()
	_, err = sigs.VerifyRequest(r, payload, fetchPublicKeyString)
	if err != nil {
		badRequest(w, err)
		return
	}
	activityArb, err := parsePayload(payload)
	if err != nil {
		badRequest(w, err)
		return
	}
	activityIRI, err := getIRI(activityArb)
	if err != nil {
		badRequest(w, err)
		return
	}
	actorArb, err := findProp(activityArb, "actor", acceptHeaders)
	if err != nil {
		badRequest(w, err)
		return
	}
	actorIRI, err := getIRI(actorArb)
	if err != nil {
		badRequest(w, err)
		return
	}
	objectArb, err := findProp(activityArb, "object", acceptHeaders)
	if err != nil {
		badRequest(w, err)
		return
	}
	objectIRI, err := getIRI(objectArb)
	if err != nil {
		badRequest(w, err)
		return
	}
	activityType, err := getType(activityArb)
	if err != nil {
		badRequest(w, err)
		return
	}
	switch activityType {
	case "Create":
		_, err = createInboxActivity(activityArb, objectArb, actorIRI.String(), recipient)
		if err != nil {
			internalServerError(w, err)
			return
		}
	case "Follow":
		if objectIRI.String() != recipient {
			badRequest(w, errors.New("wrong inbox"))
			return
		}
		_, err = createInboxReferenceActivity(activityArb, recipient, actorIRI.String(), recipient)
		if err != nil {
			internalServerError(w, err)
			return
		}
		responseArb, err := newActivityArbReference(activityIRI.String(), "Accept")
		if err != nil {
			internalServerError(w, err)
			return
		}
		responseArb["actor"] = recipient
		responseArb, err = createOutboxReferenceActivity(responseArb)
		if err != nil {
			internalServerError(w, err)
			return
		}
		fedChan <- Federation{Name: name, Recipient: actorIRI.String(), Activity: responseArb}
	case "Undo", "Accept":
		_, err = createInboxReferenceActivity(activityArb, objectIRI.String(), actorIRI.String(), recipient)
		if err != nil {
			internalServerError(w, err)
			return
		}
	default:
		badRequest(w, errors.New("unsupported activity type"))
		return
	}

	accepted(w)
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
	actor := fmt.Sprintf("%s://%s/%s/%s", config.Protocol, config.ServerName, config.Endpoints.Users, claims.Username)
	err := checkContentType(r.Header)
	if err != nil {
		badRequest(w, err)
		return
	}
	var buf bytes.Buffer
	limiter := io.LimitReader(r.Body, 1*1024*1024)
	io.Copy(&buf, limiter)
	payload := buf.Bytes()
	activityArb, err := parsePayload(payload)
	if err != nil {
		badRequest(w, err)
		return
	}
	activityType, err := getType(activityArb)
	if err != nil {
		badRequest(w, err)
		return
	}
	activityArb["actor"] = actor
	switch activityType {
	case "Create":
		objectArb, err := activityArb.GetArb("object")
		if err != nil {
			badRequest(w, err)
			return
		}
		objectArb["attributedTo"] = actor
		activityArb, err = createOutboxActivityDetail(activityArb, objectArb)
		if err != nil {
			internalServerError(w, err)
			return
		}
	case "Like":
		activityArb, err = createOutboxReferenceActivity(activityArb)
		if err != nil {
			internalServerError(w, err)
			return
		}
	case "Follow":
		activityArb, err = createOutboxReferenceActivity(activityArb)
		if err != nil {
			internalServerError(w, err)
			return
		}
		objectIRI, err := activityArb.GetURL("object")
		if err != nil {
			internalServerError(w, err)
			return
		}
		// check if the recipient is internal
		if objectIRI.Host == config.ServerName {
			// if so, generate and federate an accept
			activityIRI, err := getIRI(activityArb)
			if err != nil {
				internalServerError(w, err)
				return
			}
			responseArb, err := newActivityArbReference(activityIRI.String(), "Accept")
			if err != nil {
				internalServerError(w, err)
				return
			}
			responseArb["actor"] = objectIRI.String()
			responseArb, err = createOutboxReferenceActivity(responseArb)
			if err != nil {
				internalServerError(w, err)
				return
			}
			fedChan <- Federation{Name: name, Recipient: actor, Activity: responseArb}
		}
	case "Undo":
		activityArb, err = createOutboxReferenceActivity(activityArb)
		if err != nil {
			internalServerError(w, err)
			return
		}
	default:
		badRequest(w, errors.New("unsupported activity type"))
		return
		// Activity type is something else, save object reference (if new), Activity, and Activity_to
	}

	// activityIRI, err := getIRI(activityArb)
	// if err != nil {
	// 	internalServerError(w, err)
	// 	return
	// }

	// Get recipients
	recipients, err := getRecipients(activityArb, "to")
	if err != nil {
		log.Println(err)
	}
	// Deliver to recipients
	for _, recipient := range recipients {
		// err = addActivityTo(activityIRI.String(), recipient.String())
		// if err != nil {
		// 	log.Println(err.Error())
		// }
		fedChan <- Federation{Name: name, Recipient: recipient.String(), Activity: activityArb}
	}

	w.Header().Set("Content-Type", contentType)
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
		totalItems, err := queryFollowingTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		following := generateOrderedCollection(user.Name, config.Endpoints.Following, totalItems)
		w.Header().Set("Content-Type", contentType)
		json.NewEncoder(w).Encode(following)
		return
	}

	following, err := queryFollowingByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

	orderedItems := make([]interface{}, len(following))
	for i, actor := range following {
		orderedItems[i] = actor
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
		totalItems, err := queryFollowersTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		followers := generateOrderedCollection(user.Name, config.Endpoints.Followers, totalItems)
		w.Header().Set("Content-Type", contentType)
		json.NewEncoder(w).Encode(followers)
		return
	}

	followers, err := queryFollowersByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

	orderedItems := make([]interface{}, len(followers))
	for i, actor := range followers {
		orderedItems[i] = actor
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
		totalItems, err := queryLikedTotalItemsByUserName(user.Name)
		if err != nil {
			internalServerError(w, err)
			return
		}

		liked := generateOrderedCollection(user.Name, config.Endpoints.Liked, totalItems)
		w.Header().Set("Content-Type", contentType)
		json.NewEncoder(w).Encode(liked)
		return
	}

	liked, err := queryLikedByUserName(user.Name)
	if err != nil {
		internalServerError(w, err)
		return
	}

	orderedItems := make([]interface{}, len(liked))
	for i, activity := range liked {
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
