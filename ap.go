package main

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/cheebz/arb"
)

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

func findObject(a arb.Arb) (arb.Arb, error) {
	iri, err := a.GetURL("object")
	if err != nil {
		return a.GetArb("object")
	}
	client := http.DefaultClient
	req, err := http.NewRequest("GET", iri.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
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

// TODO: Write a method to find the Arb here instead of in the lib
