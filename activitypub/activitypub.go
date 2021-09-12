package activitypub

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/cheebz/arb"
)

var Accept = "application/activity+json"
var AcceptHeaders = http.Header{
	"Accept": []string{
		"application/activity+json",
		"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
	},
}

var ContentType = "application/activity+json"
var ContentTypeHeaders = http.Header{
	"Content-Type": []string{
		"application/activity+json",
		"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
	},
}

func CheckContentType(headers http.Header) error {
	h := headers.Values("Content-Type")
	for _, v := range h {
		log.Println("Request contains Content-Type header: " + v)
		for _, item := range ContentTypeHeaders["Content-Type"] {
			if strings.Contains(v, item) {
				return nil
			}
		}
	}
	return errors.New("invalid content-type headers")
}

func CheckAccept(headers http.Header) error {
	h := headers.Values("Accept")
	for _, v := range h {
		log.Println("Request contains Accept header: " + v)
		for _, item := range AcceptHeaders["Accept"] {
			if strings.Contains(v, item) {
				return nil
			}
		}
	}
	return errors.New("invalid accept headers")
}

var ActivityTypes = []string{"Accept", "Add", "Announce", "Arrive", "Block", "Create", "Delete", "Dislike", "Flag", "Follow", "Ignore", "Invite", "Join", "Leave", "Like", "Listen", "Move", "Offer", "Question", "Reject", "Read", "Remove", "TentativeReject", "TentativeAccept", "Travel", "Undo", "Update", "View"}
var ActorTypes = []string{"Application", "Group", "Organization", "Person", "Service"}
var ObjectTypes = []string{"Article", "Audio", "Document", "Event", "Image", "Note", "Page", "Place", "Profile", "Relationship", "Tombstone", "Video"}
var LinkTypes = []string{"Mention"}
var Audiences = []string{"to", "bto", "cc", "bcc", "audience"}

func IsActivity(t string) bool {
	for _, a := range ActivityTypes {
		if a == t {
			return true
		}
	}
	return false
}

func IsActor(t string) bool {
	for _, a := range ActorTypes {
		if a == t {
			return true
		}
	}
	return false
}

func IsObject(t string) bool {
	for _, a := range ObjectTypes {
		if a == t {
			return true
		}
	}
	return false
}

func IsLink(t string) bool {
	for _, a := range LinkTypes {
		if a == t {
			return true
		}
	}
	return false
}

func GetType(a arb.Arb) (string, error) {
	if t, err := a.GetString("type"); err == nil {
		return t, nil
	}
	if t, err := a.GetString("@type"); err == nil {
		return t, nil
	}
	return "", errors.New("unable to get type")
}

func GetIRI(a arb.Arb) (*url.URL, error) {
	if iri, err := a.GetURL("id"); err == nil {
		return iri, nil
	}
	if iri, err := a.GetURL("@id"); err == nil {
		return iri, nil
	}
	return nil, errors.New("unable to get iri")
}

func Find(iri string, headers http.Header) (arb.Arb, error) {
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

func FindProp(a arb.Arb, prop string, headers http.Header) (arb.Arb, error) {
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

func CheckContext(payload arb.Arb) error {
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

func ParsePayload(payload []byte) (arb.Arb, error) {
	payloadArb, err := arb.ReadBytes(payload)
	if err != nil {
		return nil, err
	}
	err = CheckContext(payloadArb)
	if err != nil {
		return nil, err
	}
	payloadType, err := GetType(payloadArb)
	if err != nil {
		return nil, err
	}
	var activityArb arb.Arb
	if IsObject(payloadType) {
		activityArb, err = NewActivityArb(payloadArb, "Create")
		if err != nil {
			return nil, err
		}
	}
	if IsActivity(payloadType) {
		activityArb = payloadArb
		err = FormatRecipients(activityArb)
		if err != nil {
			return nil, err
		}
	}
	if activityArb == nil {
		return nil, err
	}
	return activityArb, nil
}

func NewActivityArb(object arb.Arb, typ string) (arb.Arb, error) {
	activity := arb.New()
	activity["@context"] = []string{"https://www.w3.org/ns/activitystreams"}
	activity["type"] = typ
	err := object.PropToArray("@context")
	if err != nil {
		return nil, err
	}
	err = FormatRecipients(object)
	if err != nil {
		return nil, err
	}
	activity["object"] = object
	for _, a := range Audiences {
		if object.Exists(a) {
			activity[a] = object[a]
		}
	}
	return activity, nil
}

func NewActivityArbReference(objectIRI string, typ string) (arb.Arb, error) {
	activity := arb.New()
	activity["@context"] = []string{"https://www.w3.org/ns/activitystreams"}
	activity["type"] = typ
	activity["object"] = objectIRI
	return activity, nil
}

func FormatRecipients(a arb.Arb) error {
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

func GetRecipients(a arb.Arb, prop string) ([]*url.URL, error) {
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

func FetchPublicKeyString(keyId string) (string, error) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", keyId, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", Accept)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	key, err := arb.Read(resp.Body)
	if err != nil {
		return "", err
	}
	publicKey, err := key.GetArb("publicKey")
	if err != nil {
		return "", err
	}
	publicKeyString, err := publicKey.GetString("publicKeyPem")
	if err != nil {
		return "", err
	}
	return publicKeyString, nil
}
