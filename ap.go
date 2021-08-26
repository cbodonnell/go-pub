package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cheebz/arb"
	"github.com/cheebz/sigs"
)

var accept = "application/activity+json"
var acceptHeaders = http.Header{
	"Accept": []string{
		"application/activity+json",
		"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
	},
}

var contentType = "application/activity+json"
var contentTypeHeaders = http.Header{
	"Content-Type": []string{
		"application/activity+json",
		"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
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

var activityTypes = []string{"Accept", "Add", "Announce", "Arrive", "Block", "Create", "Delete", "Dislike", "Flag", "Follow", "Ignore", "Invite", "Join", "Leave", "Like", "Listen", "Move", "Offer", "Question", "Reject", "Read", "Remove", "TentativeReject", "TentativeAccept", "Travel", "Undo", "Update", "View"}
var actorTypes = []string{"Application", "Group", "Organization", "Person", "Service"}
var objectTypes = []string{"Article", "Audio", "Document", "Event", "Image", "Note", "Page", "Place", "Profile", "Relationship", "Tombstone", "Video"}
var linkTypes = []string{"Mention"}
var audiences = []string{"to", "bto", "cc", "bcc", "audience"}

func isActivity(t string) bool {
	for _, a := range activityTypes {
		if a == t {
			return true
		}
	}
	return false
}

func isActor(t string) bool {
	for _, a := range actorTypes {
		if a == t {
			return true
		}
	}
	return false
}

func isObject(t string) bool {
	for _, a := range objectTypes {
		if a == t {
			return true
		}
	}
	return false
}

func isLink(t string) bool {
	for _, a := range linkTypes {
		if a == t {
			return true
		}
	}
	return false
}

func getType(a arb.Arb) (string, error) {
	if t, err := a.GetString("type"); err == nil {
		return t, nil
	}
	if t, err := a.GetString("@type"); err == nil {
		return t, nil
	}
	return "", errors.New("unable to get type")
}

func getIRI(a arb.Arb) (*url.URL, error) {
	if iri, err := a.GetURL("id"); err == nil {
		return iri, nil
	}
	if iri, err := a.GetURL("@id"); err == nil {
		return iri, nil
	}
	return nil, errors.New("unable to get iri")
}

func find(iri string, headers http.Header) (arb.Arb, error) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", iri, nil)
	if err != nil {
		return nil, err
	}
	for k, l := range headers {
		for _, v := range l {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	arb, err := arb.Read(resp.Body)
	if err != nil {
		return nil, err
	}
	return arb, nil
}

func findProp(a arb.Arb, prop string, headers http.Header) (arb.Arb, error) {
	iri, err := a.GetURL(prop)
	if err != nil {
		return a.GetArb(prop)
	}
	client := http.DefaultClient
	req, err := http.NewRequest("GET", iri.String(), nil)
	if err != nil {
		return nil, err
	}
	for k, l := range headers {
		for _, v := range l {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	arb, err := arb.Read(resp.Body)
	if err != nil {
		return nil, err
	}
	return arb, nil
}

func checkContext(payload arb.Arb) error {
	err := payload.PropToArray("@context")
	if err != nil {
		return err
	}
	context, err := payload.GetArray("@context")
	if err != nil {
		return err
	}
	for _, item := range context {
		if s, ok := item.(string); ok {
			if s == "https://www.w3.org/ns/activitystreams" {
				return nil
			}
		}
	}
	return errors.New("\"https://www.w3.org/ns/activitystreams\" not in context")
}

func parsePayload(r *http.Request) (arb.Arb, error) {
	payloadArb, err := arb.Read(r.Body)
	if err != nil {
		return nil, err
	}
	err = checkContext(payloadArb)
	if err != nil {
		return nil, err
	}
	payloadType, err := payloadArb.GetString("type")
	if err != nil {
		return nil, err
	}
	var activityArb arb.Arb
	if isObject(payloadType) {
		activityArb, err = newActivityArb(payloadArb, "Create")
		if err != nil {
			return nil, err
		}
	}
	if isActivity(payloadType) {
		activityArb = payloadArb
		err = formatRecipients(activityArb)
		if err != nil {
			return nil, err
		}
	}
	if activityArb == nil {
		return nil, err
	}
	return activityArb, nil
}

func newActivityArb(object arb.Arb, typ string) (arb.Arb, error) {
	activity := arb.New()
	activity["@context"] = []string{"https://www.w3.org/ns/activitystreams"}
	activity["type"] = typ
	err := object.PropToArray("@context")
	if err != nil {
		return nil, err
	}
	err = formatRecipients(object)
	if err != nil {
		return nil, err
	}
	activity["object"] = object
	for _, a := range audiences {
		if object.Exists(a) {
			activity[a] = object[a]
		}
	}
	return activity, nil
}

func newActivityArbReference(objectIRI string, typ string) (arb.Arb, error) {
	activity := arb.New()
	activity["@context"] = []string{"https://www.w3.org/ns/activitystreams"}
	activity["type"] = typ
	activity["object"] = objectIRI
	return activity, nil
}

func formatRecipients(a arb.Arb) error {
	if a.Exists("to") {
		err := a.PropToArray("to")
		if err != nil {
			return err
		}
	}
	if a.Exists("bto") {
		err := a.PropToArray("bto")
		if err != nil {
			return err
		}
	}
	if a.Exists("cc") {
		err := a.PropToArray("cc")
		if err != nil {
			return err
		}
	}
	if a.Exists("bcc") {
		err := a.PropToArray("bcc")
		if err != nil {
			return err
		}
	}
	if a.Exists("audience") {
		err := a.PropToArray("audience")
		if err != nil {
			return err
		}
	}
	return nil
}

func getRecipients(a arb.Arb, prop string) ([]*url.URL, error) {
	urls := make([]*url.URL, 0)
	if !a.IsArray(prop) {
		return urls, nil
	}
	recipients, err := a.GetArray(prop)
	if err != nil {
		return urls, err
	}
	for _, recipient := range recipients {
		if iri, ok := recipient.(string); ok {
			if iriURL, err := url.Parse(iri); err == nil {
				urls = append(urls, iriURL)
			}
		}
	}
	return urls, nil
}

func (fed Federation) Federate() {
	logChan <- fmt.Sprintf("Federating to %s", fed.Recipient)
	recipient, err := find(fed.Recipient, acceptHeaders)
	if err != nil {
		logChan <- err.Error()
		return
	}

	recipientType, err := getType(recipient)
	if err != nil {
		logChan <- err.Error()
		return
	}

	switch recipientType {
	case "Person":
		inbox, err := recipient.GetString("inbox")
		if err != nil {
			logChan <- err.Error()
			return
		}
		fed.Post(inbox)
	case "Collection":
	case "CollectionPage":
	case "OrderedCollection":
	case "OrderedCollectionPage":
		var items []string
		orderedItems, err := recipient.GetArray("orderedItems")
		if err != nil {
			// no ordered items, get first or next
			first, err := recipient.GetString("first")
			if err != nil {
				next, err := recipient.GetString("next")
				if err != nil {
					logChan <- fmt.Sprintf("unable to federate to: %s", fed.Recipient)
				}
				fed.Recipient = next
				fed.Federate()
				return
			}
			fed.Recipient = first
			fed.Federate()
			return
		}
		for _, item := range orderedItems {
			if iri, ok := item.(string); ok {
				if iriURL, err := url.Parse(iri); err == nil {
					items = append(items, iriURL.String())
				}
			}
		}
		for _, item := range items {
			fed.Recipient = item
			fed.Federate()
		}
		// check if 'orderedItems'
		// find 'first'
		// if first get 'orderedItems'
		// find 'next' while able to
		// for each get 'orderedItems'

		// for _, recipient in orderedItems:
		// set fed.Recipient and Federate!
		return
	default:
		logChan <- fmt.Sprintf("invalid recipient type: %s", recipientType)
		return
	}
}

func (fed Federation) Post(inbox string) {
	req, err := http.NewRequest("POST", inbox, bytes.NewBuffer(fed.Data))
	if err != nil {
		logChan <- err.Error()
		return
	}
	req.Header.Add("Content-Type", contentType)

	keyID := fmt.Sprintf("%s://%s/%s/%s#main-key", config.Protocol, config.ServerName, config.Endpoints.Users, fed.Name)
	err = sigs.SignRequest(req, fed.Data, config.RSAPrivateKey, keyID)
	if err != nil {
		logChan <- err.Error()
		return
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logChan <- err.Error()
		return
	}
	defer response.Body.Close()

	logChan <- fmt.Sprintf("POST to %s", req.URL.Hostname()+req.URL.RequestURI())
	logChan <- fmt.Sprintf("%s code: %s", req.URL.Hostname()+req.URL.RequestURI(), response.Status)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logChan <- err.Error()
		return
	}
	logChan <- fmt.Sprintf("%s body: %s", req.URL.Hostname()+req.URL.RequestURI(), string(body))
}
