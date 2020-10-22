package arangodb

import driver "github.com/arangodb/go-driver"

var ARANGODB = "ARANGODB"

type database struct {
	Session       driver.Database
	ClientSession driver.Client
	Name          string
	Type          string
}

type table struct {
	Db         *database
	Name       string
	Collection driver.Collection
}

type graph struct {
	Db    *database
	Name  string
	Graph driver.Graph
}

type relation struct {
	Graph      *graph
	Name       string
	Collection driver.Collection
}

// Document contains all meta data used to identifier a document.
type Document struct {
	Key string `json:"_key,omitempty"`
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
}

type Edge struct {
	Key  string `json:"_key,omitempty"`
	ID   string `json:"_id,omitempty"`
	Rev  string `json:"_rev,omitempty"`
	From string `json:"_from,omitempty"`
	To   string `json:"_to,omitempty"`
}
