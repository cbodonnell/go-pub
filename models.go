package main

import "github.com/dgrijalva/jwt-go"

// Configuration struct
type Configuration struct {
	Debug      bool       `json:"debug"`
	Port       int        `json:"port"`
	ServerName string     `json:"serverName"`
	UserSep    string     `json:"userSep"`
	SSLCert    string     `json:"sslCert"`
	SSLKey     string     `json:"sslKey"`
	Db         DataSource `json:"db"`
	JWTKey     string     `json:"jwtKey"`
}

// Group struct
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// JWTClaims struct
type JWTClaims struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Groups   []Group `json:"groups"`
	jwt.StandardClaims
}

// DataSource struct
type DataSource struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Discoverable bool   `json:"discoverable"`
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

// Object struct
type Object struct {
	Context []string `json:"@context"`
	Id      string   `json:"id"`
	Type    string   `json:"type"`
	Name    string   `json:"name"`
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

// ActivityCollectionPage struct
type ActivityCollectionPage struct {
	Object
	PartOf       string     `json:"partOf"`
	OrderedItems []Activity `json:"orderedItems"`
}

// StringCollectionPage struct
type StringCollectionPage struct {
	Object
	PartOf       string   `json:"partOf"`
	OrderedItems []string `json:"orderedItems"`
}

// Activity struct
type Activity struct {
	Object
	Actor       string   `json:"actor"`
	To          []string `json:"to"`
	ChildObject Audio    `json:"object"`
}

// Link struct
type Link struct {
	Object
	Href      string `json:"href"`
	MediaType string `json:"mediaType"`
}

// Audio struct
type Audio struct {
	Object
	Url Link `json:"url"`
}
