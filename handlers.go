package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

	// actor, err := getActorByName(name)
	actor := Actor{
		Object: Object{
			Context: []string{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Id:   fmt.Sprintf("https://%s/%s/%s", config.ServerName, config.UserSep, name),
			Type: "Person",
		},
		Inbox:     fmt.Sprintf("https://%s/%s/%s/inbox", config.ServerName, config.UserSep, name),
		Outbox:    fmt.Sprintf("https://%s/%s/%s/outbox", config.ServerName, config.UserSep, name),
		Following: fmt.Sprintf("https://%s/%s/%s/following", config.ServerName, config.UserSep, name),
		Followers: fmt.Sprintf("https://%s/%s/%s/followers", config.ServerName, config.UserSep, name),
		Liked:     fmt.Sprintf("https://%s/%s/%s/liked", config.ServerName, config.UserSep, name),
	}

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
		inbox := generateInbox(user.Name)
		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(inbox)
	} else {
		inboxPage := generateInboxPage(name)
		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(inboxPage)
	}
}

func getOutbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	page := r.FormValue("page")
	if page != "true" {
		// outbox, err := getOutboxByName(name)
		outbox := OrderedCollection{
			Object: Object{
				Context: []string{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				Id:   fmt.Sprintf("https://%s/%s/%s/outbox", config.ServerName, config.UserSep, name),
				Type: "OrderedCollection",
			},
			TotalItems: 1,
			First:      fmt.Sprintf("https://%s/%s/%s/outbox?page=true", config.ServerName, config.UserSep, name),
			Last:       fmt.Sprintf("https://%s/%s/%s/outbox?min_id=0&page=true", config.ServerName, config.UserSep, name),
		}

		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(outbox)
	} else {
		// outbox, err := getOutboxPageByName(name)
		outboxPage := ActivityCollectionPage{
			Object: Object{
				Context: []string{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				Id:   fmt.Sprintf("https://%s/%s/%s/outbox?page=true", config.ServerName, config.UserSep, name),
				Type: "OrderedCollectionPage",
			},
			PartOf: fmt.Sprintf("https://%s/%s/%s/outbox", config.ServerName, config.UserSep, name),
			OrderedItems: append(
				[]Activity{},
				Activity{
					Object: Object{
						Context: []string{
							"https://www.w3.org/ns/activitystreams",
							"https://w3id.org/security/v1",
						},
						Type: "Create",
						Id:   fmt.Sprintf("https://%s/%s/%s/activity/1", config.ServerName, config.UserSep, name),
					},
					Actor: fmt.Sprintf("https://%s/%s/%s", config.ServerName, config.UserSep, name),
					To: []string{
						"https://www.w3.org/ns/activitystreams#Public",
					},
					ChildObject: Audio{
						Object: Object{
							Context: []string{
								"https://www.w3.org/ns/activitystreams",
								"https://w3id.org/security/v1",
							},
							Type: "Audio",
							Id:   fmt.Sprintf("https://%s/%s/%s/audio/1", config.ServerName, config.UserSep, name),
							Name: "An Audio object",
						},
						Url: Link{
							Object: Object{
								Context: []string{
									"https://www.w3.org/ns/activitystreams",
									"https://w3id.org/security/v1",
								},
								Type: "Link",
								Id:   fmt.Sprintf("https://%s/%s/%s/link/1", config.ServerName, config.UserSep, name),
								Name: "A Link object",
							},
							Href:      "https://example.org/audio.mp3",
							MediaType: "audio/mp3",
						},
					},
				},
			),
		}

		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(outboxPage)
	}

}

func getFollowing(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	page := r.FormValue("page")
	if page != "true" {
		// following, err := getFollowingByName(name)
		following := OrderedCollection{
			Object: Object{
				Context: []string{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				Id:   fmt.Sprintf("https://%s/%s/%s/following", config.ServerName, config.UserSep, name),
				Type: "OrderedCollection",
			},
			TotalItems: 1,
			First:      fmt.Sprintf("https://%s/%s/%s/following?page=true", config.ServerName, config.UserSep, name),
			Last:       fmt.Sprintf("https://%s/%s/%s/following?min_id=0&page=true", config.ServerName, config.UserSep, name),
		}

		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(following)
	} else {
		// followingPage, err := getFollowingPageByName(name)
		followingPage := StringCollectionPage{
			Object: Object{
				Context: []string{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				Id:   fmt.Sprintf("https://%s/%s/%s/following?page=true", config.ServerName, config.UserSep, name),
				Type: "OrderedCollectionPage",
			},
			PartOf: fmt.Sprintf("https://%s/%s/%s/following", config.ServerName, config.UserSep, name),
			OrderedItems: append(
				[]string{},
				fmt.Sprintf("https://%s/%s/other", config.ServerName, config.UserSep),
			),
		}

		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(followingPage)
	}
}

func getFollowers(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	page := r.FormValue("page")
	if page != "true" {
		// followers, err := getFollowersByName(name)
		followers := OrderedCollection{
			Object: Object{
				Context: []string{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				Id:   fmt.Sprintf("https://%s/%s/%s/followers", config.ServerName, config.UserSep, name),
				Type: "OrderedCollection",
			},
			TotalItems: 1,
			First:      fmt.Sprintf("https://%s/%s/%s/followers?page=true", config.ServerName, config.UserSep, name),
			Last:       fmt.Sprintf("https://%s/%s/%s/followers?min_id=0&page=true", config.ServerName, config.UserSep, name),
		}

		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(followers)
	} else {
		// followersPage, err := getFollowersPageByName(name)
		followersPage := StringCollectionPage{
			Object: Object{
				Context: []string{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				Id:   fmt.Sprintf("https://%s/%s/%s/followers?page=true", config.ServerName, config.UserSep, name),
				Type: "OrderedCollectionPage",
			},
			PartOf: fmt.Sprintf("https://%s/%s/%s/followers", config.ServerName, config.UserSep, name),
			OrderedItems: append(
				[]string{},
				fmt.Sprintf("https://%s/%s/other", config.ServerName, config.UserSep),
			),
		}

		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(followersPage)
	}
}

func getLiked(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	page := r.FormValue("page")
	if page != "true" {
		// liked, err := getLikedByName(name)
		liked := OrderedCollection{
			Object: Object{
				Context: []string{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				Id:   fmt.Sprintf("https://%s/%s/%s/liked", config.ServerName, config.UserSep, name),
				Type: "OrderedCollection",
			},
			TotalItems: 1,
			First:      fmt.Sprintf("https://%s/%s/%s/liked?page=true", config.ServerName, config.UserSep, name),
			Last:       fmt.Sprintf("https://%s/%s/%s/liked?min_id=0&page=true", config.ServerName, config.UserSep, name),
		}

		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(liked)
	} else {
		// likedPage, err := getLikedPageByName(name)
		likedPage := StringCollectionPage{
			Object: Object{
				Context: []string{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				Id:   fmt.Sprintf("https://%s/%s/%s/liked?page=true", config.ServerName, config.UserSep, name),
				Type: "OrderedCollectionPage",
			},
			PartOf: fmt.Sprintf("https://%s/%s/%s/liked", config.ServerName, config.UserSep, name),
			OrderedItems: append(
				[]string{},
				fmt.Sprintf("https://%s/%s/other/activity/1", config.ServerName, config.UserSep),
			),
		}

		w.Header().Set("Content-Type", "application/jrd+json")
		json.NewEncoder(w).Encode(likedPage)
	}
}
