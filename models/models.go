package models

import (
	"github.com/cheebz/arb"
)

// Configuration struct
type Configuration struct {
	Debug         bool       `json:"debug"`
	Port          int        `json:"port"`
	LogFile       string     `json:"logFile"`
	Protocol      string     `json:"protocol"`
	ServerName    string     `json:"serverName"`
	Auth          string     `json:"auth"`
	Client        string     `json:"client"`
	Endpoints     Endpoints  `json:"endpoints"`
	SSLCert       string     `json:"sslCert"`
	SSLKey        string     `json:"sslKey"`
	Db            DataSource `json:"db"`
	JWTKey        string     `json:"jwtKey"`
	RSAPublicKey  string     `json:"rsaPublicKey"`
	RSAPrivateKey string     `json:"rsaPrivateKey"`
}

// DataSource struct
type Endpoints struct {
	Users      string `json:"users"`
	Activities string `json:"activities"`
	Objects    string `json:"objects"`
	Inbox      string `json:"inbox"`
	Outbox     string `json:"outbox"`
	Following  string `json:"following"`
	Followers  string `json:"followers"`
	Liked      string `json:"liked"`
}

// DataSource struct
type DataSource struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

// User struct
type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Discoverable bool   `json:"discoverable"`
	IRI          string `json:"url"`
}

// Group struct
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// WebFinger struct
type WebFinger struct {
	Subject string          `json:"subject"`
	Aliases []string        `json:"aliases"`
	Links   []WebFingerLink `json:"links"`
}

// Link struct
type WebFingerLink struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

// Object struct (see: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-object)
// TODO: Should some of these be sql.NullString?
type Object struct {
	Context interface{} `json:"@context"`
	Id      string      `json:"id"`
	Type    string      `json:"type"`

	Attachment   string      `json:"attachment,omitempty"`
	AttributedTo interface{} `json:"attributedTo,omitempty"`
	Audience     []string    `json:"audience,omitempty"`
	Content      interface{} `json:"content,omitempty"`
	Name         interface{} `json:"name,omitempty"`
	EndTime      string      `json:"endTime,omitempty"`
	Generator    string      `json:"generator,omitempty"`
	Icon         string      `json:"icon,omitempty"`
	Image        string      `json:"image,omitempty"`
	InReplyTo    interface{} `json:"inReplyTo,omitempty"`
	Location     string      `json:"location,omitempty"`
	Preview      string      `json:"preview,omitempty"`
	Published    string      `json:"published,omitempty"`
	Replies      string      `json:"replies,omitempty"`
	StartTime    string      `json:"startTime,omitempty"`
	Summary      string      `json:"summary,omitempty"`
	Tag          string      `json:"tag,omitempty"`
	Updated      string      `json:"updated,omitempty"`
	Url          interface{} `json:"url,omitempty"`
	To           []string    `json:"to,omitempty"`
	Bto          []string    `json:"bto,omitempty"`
	Cc           []string    `json:"cc,omitempty"`
	Bcc          []string    `json:"bcc,omitempty"`
	MediaType    string      `json:"mediaType,omitempty"`
	Duration     string      `json:"duration,omitempty"`
}

func NewObject() Object {
	var object Object
	object.Context = []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
	}
	return object
}

// Link struct (see: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-link)
type Link struct {
	Context interface{} `json:"@context"`
	Id      string      `json:"id"`
	Type    string      `json:"type"`

	Href      string `json:"href,omitempty"`
	Rel       string `json:"rel,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
	Name      string `json:"name,omitempty"`
	HrefLang  string `json:"hreflang,omitempty"`
	Height    string `json:"height,omitempty"`
	Width     string `json:"width,omitempty"`
	Preview   string `json:"preview,omitempty"`
}

func NewLink() Link {
	var link Link
	link.Context = []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
	}
	return link
}

// Actor struct
type Actor struct {
	Object
	Inbox                     string    `json:"inbox"`
	Outbox                    string    `json:"outbox"`
	Following                 string    `json:"following,omitempty"`
	Followers                 string    `json:"followers,omitempty"`
	Liked                     string    `json:"liked,omitempty"`
	PreferredUsername         string    `json:"preferredUsername,omitempty"`
	ManuallyApprovesFollowers bool      `json:"manuallyApprovesFollowers"`
	PublicKey                 PublicKey `json:"publicKey"`
}

// PublicKey struct
type PublicKey struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

// OrderedCollection struct (see: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollection)
type OrderedCollection struct {
	Object
	TotalItems int    `json:"totalItems"`
	First      string `json:"first"`
	Last       string `json:"last"`
}

// OrderedCollectionPage struct (see: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollectionpage)
type OrderedCollectionPage struct {
	Object
	PartOf       string        `json:"partOf"`
	OrderedItems []interface{} `json:"orderedItems"`
	Next         string        `json:"next,omitempty"`
	Prev         string        `json:"prev,omitempty"`
}

// PostActivityResource struct
type PostActivityResource struct {
	Object
	Actor       string `json:"actor"`
	ChildObject Object `json:"object"`
}

// Activity struct (see: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-activity)
type Activity struct {
	Object
	Actor       string      `json:"actor"`
	ChildObject interface{} `json:"object"`
}

func NewActivity() Activity {
	var activity Activity
	activity.Context = []interface{}{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
	}
	return activity
}

type Federation struct {
	Name      string
	Recipient string
	Activity  arb.Arb
}

type CheckResponse struct {
	Exists      bool   `json:"exists"`
	ActivityIRI string `json:"iri"`
}
