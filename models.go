package main

// Configuration struct
type Configuration struct {
	Debug      bool       `json:"debug"`
	Port       int        `json:"port"`
	Protocol   string     `json:"protocol"`
	ServerName string     `json:"serverName"`
	Auth       string     `json:"auth"`
	Client     string     `json:"client"`
	Endpoints  Endpoints  `json:"endpoints"`
	SSLCert    string     `json:"sslCert"`
	SSLKey     string     `json:"sslKey"`
	Db         DataSource `json:"db"`
	JWTKey     string     `json:"jwtKey"`
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

// Group struct
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// HomeData struct - Data sent to the index.html template
type HomeData struct {
	Claims         *JWTClaims
	User           User
	ServerName     string
	UsersEndpoint  string
	OutboxEndpoint string
	Auth           string
}

// User struct
type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Discoverable bool   `json:"discoverable"`
	IRI          string `json:"url"`
}

// ActivityOLD struct
type ActivityOLD struct {
	ID       int      `json:"id"`
	UserName string   `json:"userName"`
	Type     string   `json:"type"`
	To       []string `json:"to"`
}

// Note struct
type Note struct {
	ID       int         `json:"id"`
	UserName string      `json:"userName"`
	Content  string      `json:"content"`
	Activity ActivityOLD `json:"activity"`
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
	Context []interface{} `json:"@context"`
	Id      string   `json:"id"`
	Type    string   `json:"type"`

	Attachment   string   `json:"attachment,omitempty"`
	AttributedTo string   `json:"attributedTo,omitempty"`
	Audience     []string `json:"audience,omitempty"`
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
	Url          interface{}    `json:"url,omitempty"`
	To           []string `json:"to,omitempty"`
	Bto          []string `json:"bto,omitempty"`
	Cc           []string `json:"cc,omitempty"`
	Bcc          []string `json:"bcc,omitempty"`
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
	PreferredUsername string `json:"preferredUsername"`
	ManuallyApprovesFollowers bool `json:"manuallyApprovesFollowers"`
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
