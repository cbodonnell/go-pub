package main

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/cheebz/arb"
)

var activityTypes = []string{"Accept", "Add", "Announce", "Arrive", "Block", "Create", "Delete", "Dislike", "Flag", "Follow", "Ignore", "Invite", "Join", "Leave", "Like", "Listen", "Move", "Offer", "Question", "Reject", "Read", "Remove", "TentativeReject", "TentativeAccept", "Travel", "Undo", "Update", "View"}
var actorTypes = []string{"Application", "Group", "Organization", "Person", "Service"}
var objectTypes = []string{"Article", "Audio", "Document", "Event", "Image", "Note", "Page", "Place", "Profile", "Relationship", "Tombstone", "Video"}
var linkTypes = []string{"Mention"}

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

func findObject(a arb.Arb, headers http.Header) (arb.Arb, error) {
	iri, err := a.GetURL("object")
	if err != nil {
		return a.GetArb("object")
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

func createActivity(object arb.Arb) (arb.Arb, error) {
	activity := arb.New()
	activity["@context"] = []string{"https://www.w3.org/ns/activitystreams"}
	activity["type"] = "Create"
	err := object.PropToArray("@context")
	if err != nil {
		return nil, err
	}
	err = formatRecipients(object)
	if err != nil {
		return nil, err
	}
	// activity["id"] = "Create" // This is auto-generated and added later
	// activity["actor"] = "username" // This is from auth
	activity["object"] = object
	activity["to"] = object["to"]
	activity["bto"] = object["bto"]
	activity["cc"] = object["cc"]
	activity["bcc"] = object["bcc"]
	activity["audience"] = object["audience"]
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
