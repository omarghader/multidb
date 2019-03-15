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
	Graph(name string, from []string, to []string) Graph
	ExecQuery(query string, params interface{}) ([]interface{}, error)
}

type Table interface {
	Create() error
	Exists() bool
	Drop() error
	Insert(string, interface{}) error
	Find(string) (interface{}, error)
	Update(string, interface{}) error
	Delete(string) error
}

type Graph interface {
	Create() error
	Relation(name string) Relation
}

type Relation interface {
	Insert(string, string, interface{}) error
	Find(string) (interface{}, error)
	Update(string, interface{}) error
	Delete(string) error
}
