package activitypub

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/cheebz/go-pub/config"
	"github.com/cheebz/go-pub/models"
	"github.com/cheebz/go-pub/repositories"
	"github.com/cheebz/sigs"
)

type Federator struct {
	conf config.Configuration
	repo repositories.Repository
}

func NewFederator(_conf config.Configuration, _repo repositories.Repository) Federator {
	return Federator{
		conf: _conf,
		repo: _repo,
	}
}

func (f *Federator) Federate(fed models.Federation) {
	log.Println(fmt.Sprintf("Federating to %s", fed.Recipient))
	recipient, err := Find(fed.Recipient, AcceptHeaders)
	if err != nil {
		log.Println(err)
		return
	}
	recipientType, err := GetType(recipient)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fmt.Sprintf("%s is of type %s", fed.Recipient, recipientType))

	switch recipientType {
	case "Person", "Service":
		activityIRI, err := GetIRI(fed.Activity)
		if err != nil {
			log.Println(err)
			return
		}
		recipientIRI, err := GetIRI(recipient)
		if err != nil {
			log.Println(err)
			return
		}
		if f.repo.ActivityToExists(activityIRI.String(), recipientIRI.String()) {
			return
		}
		err = f.repo.AddActivityTo(activityIRI.String(), recipientIRI.String())
		if err != nil {
			log.Println(err)
			return
		}
		if recipientIRI.Host != f.conf.ServerName {
			inbox, err := recipient.GetString("inbox")
			if err != nil {
				log.Println(err)
				return
			}
			f.post(fed, inbox)
			return
		}
		log.Println(fmt.Sprintf("%s is a local user", fed.Recipient))
		return
	case "Collection", "CollectionPage", "OrderedCollection", "OrderedCollectionPage":
		log.Println(fmt.Sprintf("%s is a collection", fed.Recipient))
		var items []string
		orderedItems, err := recipient.GetArray("orderedItems")
		if err != nil {
			// no ordered items, get first or next
			first, err := recipient.GetString("first")
			if err != nil {
				next, err := recipient.GetString("next")
				if err != nil {
					log.Println(fmt.Sprintf("unable to federate to: %s", fed.Recipient))
				}
				fed.Recipient = next
				f.Federate(fed)
				return
			}
			fed.Recipient = first
			f.Federate(fed)
			return
		}
		log.Println(fmt.Sprintf("retrieved orderedItems from %s", fed.Recipient))
		for _, item := range orderedItems {
			if iri, ok := item.(string); ok {
				if iriURL, err := url.Parse(iri); err == nil {
					items = append(items, iriURL.String())
				}
			}
		}
		for _, item := range items {
			fed.Recipient = item
			f.Federate(fed)
		}
		return
	default:
		log.Println(fmt.Sprintf("invalid recipient type: %s", recipientType))
		return
	}
}

func (f *Federator) post(fed models.Federation, inbox string) {
	req, err := http.NewRequest("POST", inbox, bytes.NewBuffer(fed.Activity.ToBytes()))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Content-Type", ContentType)

	keyID := fmt.Sprintf("%s://%s/%s/%s#main-key", f.conf.Protocol, f.conf.ServerName, f.conf.Endpoints.Users, fed.Name)
	err = sigs.SignRequest(req, fed.Activity.ToBytes(), f.conf.RSAPrivateKey, keyID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(fmt.Sprintf("POST to %s", req.URL.Hostname()+req.URL.RequestURI()))
	// Is it possible not to wait for this and have it also done concurrently?
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	log.Println(fmt.Sprintf("%s code: %s", req.URL.Hostname()+req.URL.RequestURI(), response.Status))
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fmt.Sprintf("%s body: %s", req.URL.Hostname()+req.URL.RequestURI(), string(body)))
}
