package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cheebz/arb"
	"github.com/cheebz/sigs"
)

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

func federate(name string, inbox string, data []byte) {
	req, err := http.NewRequest("POST", inbox, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	for k, l := range contentTypeHeaders {
		for _, v := range l {
			req.Header.Add(k, v)
		}
	}

	// TODO: Refactor into sigs
	headers := []string{"(request-target)", "date", "host", "content-type", "digest"}
	var signedLines []string
	for _, h := range headers {
		var s string
		switch h {
		case "(request-target)":
			s = strings.ToLower(req.Method) + " " + req.URL.RequestURI()
		case "date":
			s = req.Header.Get(h)
			if s == "" {
				s = time.Now().UTC().Format(http.TimeFormat)
				req.Header.Set(h, s)
			}
		case "host":
			s = req.Header.Get(h)
			if s == "" {
				s = req.URL.Hostname()
				req.Header.Set(h, s)
			}
		case "content-type":
			s = req.Header.Get(h)
		case "digest":
			s = req.Header.Get(h)
			if s == "" {
				digest, err := sigs.Digest(data)
				if err != nil {
					fmt.Println(err)
				}
				s = fmt.Sprintf("SHA-256=%s", digest)
				req.Header.Set(h, s)
			}
		}
		signedLines = append(signedLines, h+": "+s)
	}
	signedString := strings.Join(signedLines, "\n")
	// fmt.Println(signedString)

	key, err := sigs.ReadPrivateKey([]byte(config.RSAPrivateKey))
	if err != nil {
		fmt.Println(err)
	}
	sig, err := sigs.SignString(key, signedString)
	if err != nil {
		fmt.Println(err)
	}

	sigHeader := fmt.Sprintf(`keyId="%s",algorithm="%s",headers="%s",signature="%s"`,
		fmt.Sprintf("%s://%s/%s/%s#main-key", config.Protocol, config.ServerName, config.Endpoints.Users, name),
		"rsa-sha256",
		strings.Join(headers, " "),
		sig,
	)
	fmt.Println(sigHeader)
	req.Header.Set("Signature", sigHeader)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("response Body:", string(body))
}
