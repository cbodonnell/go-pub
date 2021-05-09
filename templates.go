package main

import "fmt"

func generateWebFinger(name string) WebFinger {
	return WebFinger{
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
}

func generateInbox(name string) OrderedCollection {
	return OrderedCollection{
		Object: Object{
			Context: []string{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Id:   fmt.Sprintf("https://%s/%s/%s/inbox", config.ServerName, config.UserSep, name),
			Type: "OrderedCollection",
		},
		TotalItems: 1,
		First:      fmt.Sprintf("https://%s/%s/%s/inbox?page=true", config.ServerName, config.UserSep, name),
		Last:       fmt.Sprintf("https://%s/%s/%s/inbox?min_id=0&page=true", config.ServerName, config.UserSep, name),
	}
}

func generateInboxPage(name string) ActivityCollectionPage {
	return ActivityCollectionPage{
		Object: Object{
			Context: []string{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Id:   fmt.Sprintf("https://%s/%s/%s/inbox?page=true", config.ServerName, config.UserSep, name),
			Type: "OrderedCollectionPage",
		},
		PartOf: fmt.Sprintf("https://%s/%s/%s/inbox", config.ServerName, config.UserSep, name),
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
				Actor: fmt.Sprintf("https://%s/%s/other", config.ServerName, config.UserSep),
				To: []string{
					fmt.Sprintf("https://%s/%s/%s", config.ServerName, config.UserSep, name),
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
}
