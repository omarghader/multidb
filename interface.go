package multidb

type ConnectionOptions struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type Database interface {
	Create() error
	Exists() bool
	Drop() error
	Table(name string) Table
	Graph(name string) Graph
	ExecQuery(query string, params map[string]interface{}, res []interface{}) ([]interface{}, error)
}

type Table interface {
	Create() error
	Exists() bool
	Drop() error
	CRUD
}

type Graph interface {
	Create() error
	Relation(name string) Relation
}

type Relation interface {
	CRUD
}

type CRUD interface {
	Insert(data interface{}, res interface{}) (interface{}, error)
	Find(id string, res interface{}) (interface{}, error)
	Update(id string, data interface{}, res interface{}) (interface{}, error)
	Delete(id string, res interface{}) (interface{}, error)
}
