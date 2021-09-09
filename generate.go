package main

import (
	"fmt"

	"github.com/cheebz/go-pub/config"
)

// TODO: Move this stuff into a package

func generateWebFinger(name string) WebFinger {
	return WebFinger{
		Subject: fmt.Sprintf("acct:%s@%s", name, config.C.ServerName),
		Aliases: []string{
			fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name),
		},
		Links: []WebFingerLink{
			{
				Rel:  "http://webfinger.net/rel/profile-page",
				Type: "text/html",
				Href: fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name),
			},
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name),
			},
			{
				Rel:  "http://ostatus.org/schema/1.0/subscribe",
				Href: fmt.Sprintf("%s://%s/%s", config.C.Protocol, config.C.ServerName, "/authorize_interaction?uri={uri}"),
			},
		},
	}
}

func generateActor(name string) Actor {
	return Actor{
		Object: Object{
			Context: []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
				map[string]interface{}{
					"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
				},
			},
			Id:      fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name),
			Type:    "Person",
			Name:    name,
			Url:     fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name),
			Summary: fmt.Sprintf("Summary of %s to come...", name), // TODO: Implement this
		},
		Inbox:                     fmt.Sprintf("%s://%s/%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, config.C.Endpoints.Inbox),
		Outbox:                    fmt.Sprintf("%s://%s/%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, config.C.Endpoints.Outbox),
		Following:                 fmt.Sprintf("%s://%s/%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, config.C.Endpoints.Following),
		Followers:                 fmt.Sprintf("%s://%s/%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, config.C.Endpoints.Followers),
		Liked:                     fmt.Sprintf("%s://%s/%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, config.C.Endpoints.Liked),
		PreferredUsername:         name,
		ManuallyApprovesFollowers: false, // TODO: Implement this
		PublicKey: PublicKey{
			ID:           fmt.Sprintf("%s://%s/%s/%s#main-key", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name),
			Owner:        fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name),
			PublicKeyPem: config.C.RSAPublicKey,
		},
	}
}

func generateNewActivity() Activity {
	var activity Activity
	activity.Context = []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
	}
	return activity
}

func generateNewObject() Object {
	var object Object
	object.Context = []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
	}
	return object
}

func generateOrderedCollection(name string, endpoint string, totalItems int) OrderedCollection {
	return OrderedCollection{
		Object: Object{
			Context: []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Id:   fmt.Sprintf("%s://%s/%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, endpoint),
			Type: "OrderedCollection",
		},
		TotalItems: totalItems,
		First:      fmt.Sprintf("%s://%s/%s/%s/%s?page=true", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, endpoint),
		Last:       fmt.Sprintf("%s://%s/%s/%s/%s?min_id=0&page=true", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, endpoint),
	}
}

func generateOrderedCollectionPage(name string, endpoint string, orderedItems []interface{}) OrderedCollectionPage {
	return OrderedCollectionPage{
		Object: Object{
			Context: []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Id:   fmt.Sprintf("%s://%s/%s/%s/%s?page=true", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, endpoint),
			Type: "OrderedCollectionPage",
		},
		PartOf:       fmt.Sprintf("%s://%s/%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, name, endpoint),
		OrderedItems: orderedItems,
	}
}

func generatePostActivity(post Note) PostActivityResource {
	// TODO: Get this array
	to := []string{
		"https://www.w3.org/ns/activitystreams#Public",
	}
	for _, url := range post.Activity.To {
		to = append(to, url)
	}

	return PostActivityResource{
		Object: Object{
			Context: []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Type: post.Activity.Type,
			Id:   fmt.Sprintf("%s://%s/%s/%s/activities/%d", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, post.Activity.UserName, post.Activity.ID),
			To:   to,
		},
		Actor: fmt.Sprintf("%s://%s/%s/%s", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, post.Activity.UserName),
		ChildObject: Object{
			Context: []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Type:    "Note",
			Id:      fmt.Sprintf("%s://%s/%s/%s/posts/%d", config.C.Protocol, config.C.ServerName, config.C.Endpoints.Users, post.UserName, post.ID),
			Name:    fmt.Sprintf("A note from %s", post.UserName),
			Content: post.Content,
		},
	}
}
