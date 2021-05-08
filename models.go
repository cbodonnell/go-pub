package main

import "github.com/dgrijalva/jwt-go"

// Configuration struct
type Configuration struct {
	Debug      bool   `json:"debug"`
	Port       int    `json:"port"`
	SSLCert    string `json:"sslCert"`
	SSLKey     string `json:"sslKey"`
	JWTKey     string `json:"jwtKey"`
	ServerName string `json:"serverName"`
	UserSep    string `json:"userSep"`
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

// WebFinger struct
type WebFinger struct {
	Subject string   `json:"subject"`
	Aliases []string `json:"aliases"`
	Links   []Link   `json:"links"`
}

// Link struct
type Link struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

// TODO: Can inherit @context and stuff?

// Actor struct
type Actor struct {
	Context []string `json:"@context"`
	Id      string   `json:"id"`
	Type    string   `json:"type"`
	Inbox   string   `json:"inbox"`
	Outbox  string   `json:"outbox"`
}

// Mailbox struct
type Mailbox struct {
	Context    []string `json:"@context"`
	Id         string   `json:"id"`
	Type       string   `json:"type"`
	TotalItems int      `json:"totalItems"`
	First      string   `json:"first"`
	Last       string   `json:"last"`
}

// Mailbox struct
type MailboxPage struct {
	Context      []string   `json:"@context"`
	Id           string     `json:"id"`
	Type         string     `json:"type"`
	PartOf       string     `json:"partOf"`
	OrderedItems []Activity `json:"orderedItems"`
}

type Activity struct {
	Context []string `json:"@context"`
	Id      string   `json:"id"`
	Type    string   `json:"type"`
	Actor   string   `json:"name"`
	To      []string `json:"to"`
	Object  Object   `json:"object"`
}

// TODO: Make a generic type and inherit
type Object struct {
	Context []string `json:"@context"`
	Id      string   `json:"id"`
	Type    string   `json:"type"`
	// Actor		string		`json:"actor"`
	Name      string `json:"name"`
	MediaType string `json:"mediaType"`
	Content   string `json:"content"`
}
