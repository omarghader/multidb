# multidb
The aim of this library is to make one interface for multiple databases.

# Todo :
* Relational DB
	* [ ]  MySQL
	* [ ]  PostgreSql
* NoSQL DB
	* [ ] Mongodb
	* [ ] Arangodb
* Graph DB
	* [ ] Neo4j
	* [x] Arangodb


# Example

```go
db := arangodb.NewArangodb(multidb.ConnectionOptions{
  Host:     "localhost",
  Port:     "8529",
  Username: "root",
  Password: "root",
  DBName:   "db",
})

db.Create()

// Create a table
db.Table("table_test").Create()

// Insert documents
type DocTest struct {
  Key   string `json:"_key,omitempty"`
  ID    string `json:"_id,omitempty"`
  Name  string
}

doc1 := DocTest{Key:  "doc1", Name: "document1"}
var doc1Res DocTest
db.Table("table_test").Insert(doc1, &doc1Res)

doc2 := DocTest{Key:  "doc2", Name: "document2"}
var doc2Res DocTest
db.Table("table_test").Insert(doc2, &doc2Res)


// Create a graph
db.Graph("graph_test").Create()

// Create relation
db.Graph("graph_test").Relation("is_related").Insert(
  arangodb.Edge{
    From: doc1Res.ID,
    To:   doc2Res.ID,
  }, map[string]interface{}{"prop1": "friend"})

```

# API

```go
type Database interface {
	Create() error
	Exists() bool
	Drop() error
	Table(name string) Table
	Graph(name string) Graph
	ExecQuery(query string, params map[string]interface{}, res []interface{}) ([]interface{}, error)
}

// Collection or table
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
  Insert(data interface{}, res interface{}) (interface{}, error)
	Find(id string, res interface{}) (interface{}, error)
	Update(id string, data interface{}, res interface{}) (interface{}, error)
	Delete(id string, res interface{}) (interface{}, error)
}

type CRUD interface {
	Insert(data interface{}, res interface{}) (interface{}, error)
	Find(id string, res interface{}) (interface{}, error)
	Update(id string, data interface{}, res interface{}) (interface{}, error)
	Delete(id string, res interface{}) (interface{}, error)
}
