package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func webFinger(w http.ResponseWriter, r *http.Request) {
	orig := r.FormValue("resource")

	log.Printf("resource: %s", orig)

	if strings.HasPrefix(orig, "acct:") {
		orig = orig[5:]
	}

	// webfinger := getWebFingerByAcct
	name := orig
	idx := strings.LastIndexByte(name, '/')
	if idx != -1 {
		name = name[idx+1:]
		if fmt.Sprintf("https://%s:/%s/%s", config.ServerName, config.UserSep, name) != orig {
			log.Printf("foreign request rejected")
			badRequest(w, errors.New("foreign request rejected"))
			return
		}
	} else {
		idx = strings.IndexByte(name, '@')
		if idx != -1 {
			name = name[:idx]
			if !(name+"@"+config.ServerName == orig) {
				log.Printf("foreign request rejected")
				badRequest(w, errors.New("foreign request rejected"))
				return
			}
		}
	}
	// user, err := getUserByName(name)
	// if err != nil {
	// 	http.NotFound(w, r)
	// 	return
	// }
	// if stealthmode(user.ID, r) {
	// 	http.NotFound(w, r)
	// 	return
	// }

	webfinger := WebFinger{
		Subject: fmt.Sprintf("acct:%s@%s", name, config.ServerName),
		Aliases: []string{
			fmt.Sprintf("https://%s/%s/%s", config.ServerName, config.UserSep, name),
		},
		Links: append(
			[]WebFingerLink{},
			WebFingerLink{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("https://%s/%s/%s", config.ServerName, config.UserSep, name),
			},
		),
	}

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
		Inbox:  fmt.Sprintf("https://%s/%s/%s/inbox", config.ServerName, config.UserSep, name),
		Outbox: fmt.Sprintf("https://%s/%s/%s/outbox", config.ServerName, config.UserSep, name),
	}

	w.Header().Set("Content-Type", "application/jrd+json")
	json.NewEncoder(w).Encode(actor)
}

func getInbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// inbox, err := getInboxByName(name)
	inbox := Mailbox{
		Object: Object{
			Context: []string{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Id:   fmt.Sprintf("https://%s/%s/%s/inbox", config.ServerName, config.UserSep, name),
			Type: "OrderedCollection",
		},
		TotalItems: 1,
		First:      fmt.Sprintf("https://%s/%s/%s/inbox", config.ServerName, config.UserSep, name),
		Last:       fmt.Sprintf("https://%s/%s/%s/inbox?min_id=0", config.ServerName, config.UserSep, name),
	}

	w.Header().Set("Content-Type", "application/jrd+json")
	json.NewEncoder(w).Encode(inbox)
}

func getOutbox(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	page := r.FormValue("page")
	if page != "true" {
		// outbox, err := getOutboxByName(name)
		outbox := Mailbox{
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
		outboxPage := MailboxPage{
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
