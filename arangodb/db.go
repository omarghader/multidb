package multidb

var ARANGODB = "ARANGODB"

type database struct {
	Session       interface{}
	ClientSession interface{}
	Name          string
	Type          string
}

type table struct {
	Db   *database
	Name string
}

type graph struct {
	Db       *database
	Name     string
	fromType []string
	toType   []string
}

type relation struct {
	Graph *graph
	Name  string
}
