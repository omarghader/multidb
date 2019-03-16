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
	Graph(name string, fromTables []string, toTables []string) Graph
	ExecQuery(query string, params map[string]interface{}) ([]interface{}, error)
}

type Table interface {
	Create() error
	Exists() bool
	Drop() error
	Insert(string, interface{}) (interface{}, error)
	Find(string) (interface{}, error)
	Update(string, interface{}) (interface{}, error)
	Delete(string) (interface{}, error)
}

type Graph interface {
	Create() error
	Relation(name string) Relation
}

type Relation interface {
	Create(fromNode, toNode string, params interface{}) error
	Insert(string, string, interface{}) (interface{}, error)
	Find(string) (interface{}, error)
	Update(string, interface{}) (interface{}, error)
	Delete(string) (interface{}, error)
}
