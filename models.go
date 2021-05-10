package main

import "github.com/dgrijalva/jwt-go"

// Configuration struct
type Configuration struct {
	Debug      bool       `json:"debug"`
	Port       int        `json:"port"`
	ServerName string     `json:"serverName"`
	Endpoints  Endpoints  `json:"endpoints"`
	SSLCert    string     `json:"sslCert"`
	SSLKey     string     `json:"sslKey"`
	Db         DataSource `json:"db"`
	JWTKey     string     `json:"jwtKey"`
}

// DataSource struct
type Endpoints struct {
	Users     string `json:"users"`
	Inbox     string `json:"inbox"`
	Outbox    string `json:"outbox"`
	Following string `json:"following"`
	Followers string `json:"followers"`
	Liked     string `json:"liked"`
}

// DataSource struct
type DataSource struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

// JWTClaims struct
type JWTClaims struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Groups   []Group `json:"groups"`
	jwt.StandardClaims
}

// Group struct
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// User struct
type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Discoverable bool   `json:"discoverable"`
	URL          string `json:"url"`
}

// Activity struct
type Activity struct {
	ID       int      `json:"id"`
	UserName string   `json:"userName"`
	Type     string   `json:"type"`
	To       []string `json:"to"`
}

// Post struct
type Post struct {
	ID       int      `json:"id"`
	UserName string   `json:"userName"`
	Content  string   `json:"content"`
	Activity Activity `json:"activity"`
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
type Object struct {
	Context []string `json:"@context"`
	Id      string   `json:"id"`
	Type    string   `json:"type"`

	Attachment   string   `json:"attachment,omitempty"`
	AttributedTo string   `json:"attributedTo,omitempty"`
	Audience     string   `json:"audience,omitempty"`
	Content      string   `json:"content,omitempty"`
	Name         string   `json:"name,omitempty"`
	EndTime      string   `json:"endTime,omitempty"`
	Generator    string   `json:"generator,omitempty"`
	Icon         string   `json:"icon,omitempty"`
	Image        string   `json:"image,omitempty"`
	InReplyTo    string   `json:"inReplyTo,omitempty"`
	Location     string   `json:"location,omitempty"`
	Preview      string   `json:"preview,omitempty"`
	Published    string   `json:"published,omitempty"`
	Replies      string   `json:"replies,omitempty"`
	StartTime    string   `json:"startTime,omitempty"`
	Summary      string   `json:"summary,omitempty"`
	Tag          string   `json:"tag,omitempty"`
	Updated      string   `json:"updated,omitempty"`
	Url          *Link    `json:"url,omitempty"`
	To           []string `json:"to,omitempty"`
	Bto          string   `json:"bto,omitempty"`
	Cc           string   `json:"cc,omitempty"`
	Bcc          string   `json:"bcc,omitempty"`
	MediaType    string   `json:"mediaType,omitempty"`
	Duration     string   `json:"duration,omitempty"`
}

// Link struct (see: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-link)
type Link struct {
	Context []string `json:"@context"`
	Id      string   `json:"id"`
	Type    string   `json:"type"`

	Href      string `json:"href,omitempty"`
	Rel       string `json:"rel,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
	Name      string `json:"name,omitempty"`
	HrefLang  string `json:"hreflang,omitempty"`
	Height    string `json:"height,omitempty"`
	Width     string `json:"width,omitempty"`
	Preview   string `json:"preview,omitempty"`
}

// Actor struct
type Actor struct {
	Object
	Inbox     string `json:"inbox"`
	Outbox    string `json:"outbox"`
	Following string `json:"following"`
	Followers string `json:"followers"`
	Liked     string `json:"liked"`
}

// OrderedCollection struct
type OrderedCollection struct {
	Object
	TotalItems int    `json:"totalItems"`
	First      string `json:"first"`
	Last       string `json:"last"`
}

// OrderedCollectionPage struct
type OrderedCollectionPage struct {
	Object
	PartOf       string        `json:"partOf"`
	OrderedItems []interface{} `json:"orderedItems"`
}

// PostActivityResource struct
type PostActivityResource struct {
	Object
	Actor       string `json:"actor"`
	ChildObject Object `json:"object"`
}
