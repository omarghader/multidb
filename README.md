# multidb
A dao for multiple databases. one interface multiple implementations
# multidb
The aim of this library is to make one interface for multiple databases.

# API

```go
type Dao interface {
	Connect(options) error
}

type CRUDInterface interface {
	Insert(model) (model, error)
    Find(filter) ([]model, error)
    Update(model) (model, error)
    Delete(model) (model, error)
}

```

# Todo :
* Relational DB
	* [ ]  MySQL
	* [ ]  PostgreSql

* NoSQL DB
	* [ ] Mongodb
	* [ ] Arangodb

* Graph DB
	* [ ] Neo4j
	* [ ] Arangodb

# Todo :
* [ ]  MySQL
