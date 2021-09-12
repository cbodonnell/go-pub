package resources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/models"
)

type ActivityPubResource struct {
	conf config.Configuration
}

func NewActivityPubResource(_conf config.Configuration) Resource {
	return &ActivityPubResource{
		conf: _conf,
	}
}

func (r *ActivityPubResource) ParseResource(resource string) (string, error) {
	if strings.HasPrefix(resource, "acct:") {
		resource = resource[5:]
	}
	name := resource
	idx := strings.LastIndexByte(name, '/')
	if idx != -1 {
		name = name[idx+1:]
		if fmt.Sprintf("%s/%s/%s", r.conf.ServerName, r.conf.Endpoints.Users, name) != resource {
			return name, errors.New("foreign request rejected")
		}
	} else {
		idx = strings.IndexByte(name, '@')
		if idx != -1 {
			name = name[:idx]
			if !(name+"@"+r.conf.ServerName == resource) {
				return name, errors.New("foreign request rejected")
			}
		}
	}
	return name, nil
}

func (r *ActivityPubResource) GenerateWebFinger(name string) models.WebFinger {
	return models.WebFinger{
		Subject: fmt.Sprintf("acct:%s@%s", name, r.conf.ServerName),
		Aliases: []string{
			fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
		},
		Links: []models.WebFingerLink{
			{
				Rel:  "http://webfinger.net/rel/profile-page",
				Type: "text/html",
				Href: fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
			},
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
			},
			{
				Rel:  "http://ostatus.org/schema/1.0/subscribe",
				Href: fmt.Sprintf("%s://%s/%s", r.conf.Protocol, r.conf.ServerName, "/authorize_interaction?uri={uri}"),
			},
		},
	}
}

func (r *ActivityPubResource) GenerateActor(name string) models.Actor {
	return models.Actor{
		Object: models.Object{
			Context: []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
				map[string]interface{}{
					"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
				},
			},
			Id:      fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
			Type:    "Person",
			Name:    name,
			Url:     fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
			Summary: fmt.Sprintf("Summary of %s to come...", name), // TODO: Implement this
		},
		Inbox:                     fmt.Sprintf("%s://%s/%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, r.conf.Endpoints.Inbox),
		Outbox:                    fmt.Sprintf("%s://%s/%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, r.conf.Endpoints.Outbox),
		Following:                 fmt.Sprintf("%s://%s/%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, r.conf.Endpoints.Following),
		Followers:                 fmt.Sprintf("%s://%s/%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, r.conf.Endpoints.Followers),
		Liked:                     fmt.Sprintf("%s://%s/%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, r.conf.Endpoints.Liked),
		PreferredUsername:         name,
		ManuallyApprovesFollowers: false, // TODO: Implement this
		PublicKey: models.PublicKey{
			ID:           fmt.Sprintf("%s://%s/%s/%s#main-key", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
			Owner:        fmt.Sprintf("%s://%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name),
			PublicKeyPem: r.conf.RSAPublicKey,
		},
	}
}

func (r *ActivityPubResource) GenerateOrderedCollection(name string, endpoint string, totalItems int) models.OrderedCollection {
	return models.OrderedCollection{
		Object: models.Object{
			Context: []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Id:   fmt.Sprintf("%s://%s/%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, endpoint),
			Type: "OrderedCollection",
		},
		TotalItems: totalItems,
		First:      fmt.Sprintf("%s://%s/%s/%s/%s?page=true", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, endpoint),
		Last:       fmt.Sprintf("%s://%s/%s/%s/%s?min_id=0&page=true", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, endpoint),
	}
}

func (r *ActivityPubResource) GenerateOrderedCollectionPage(name string, endpoint string, orderedItems []interface{}) models.OrderedCollectionPage {
	return models.OrderedCollectionPage{
		Object: models.Object{
			Context: []interface{}{
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
			},
			Id:   fmt.Sprintf("%s://%s/%s/%s/%s?page=true", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, endpoint),
			Type: "OrderedCollectionPage",
		},
		PartOf:       fmt.Sprintf("%s://%s/%s/%s/%s", r.conf.Protocol, r.conf.ServerName, r.conf.Endpoints.Users, name, endpoint),
		OrderedItems: orderedItems,
	}
}
