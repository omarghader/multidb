package arangodb

import (
	"errors"
	"fmt"
	"strings"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/omarghader/multidb"
	"github.com/sirupsen/logrus"
)

func NewArangodb(options multidb.ConnectionOptions) multidb.Database {
	db := &database{Name: options.DBName, Type: ARANGODB}

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{fmt.Sprintf("http://%s:%s", options.Host, options.Port)},
	})

	if err != nil {
		logrus.Errorf("%s: %s", multidb.EXCEPTION_CONNECTION_FAILED, err.Error())
		return nil
	}

	session, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(options.Username, options.Password),
	})

	if err != nil {
		logrus.Errorf("%s: %s", multidb.EXCEPTION_CONNECTION_FAILED, err.Error())
		return nil
	}

	db.ClientSession = session

	if db.Exists() {
		dbSession, err := session.Database(nil, db.Name)
		if err != nil {
			logrus.Errorf("%s: %s\n", multidb.EXCEPTION_CONNECTION_FAILED, err.Error())
			return nil
		}
		db.Session = dbSession
	}

	return db
}

func (d *database) getSession() driver.Database {
	return d.Session
}

func (d *database) getClientSession() driver.Client {
	return d.ClientSession
}

func (d *database) Create() error {
	dbSession, err := d.getClientSession().CreateDatabase(nil, d.Name, nil)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_DB_CREATE, err.Error())
		return err
	}
	d.Session = dbSession
	return nil
}

func (d *database) Exists() bool {
	exists, err := d.getClientSession().DatabaseExists(nil, d.Name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_CONNECTION_FAILED, err.Error())
		return false
	}

	return exists
}

func (d *database) Drop() error {
	db, err := d.getClientSession().Database(nil, d.Name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_DB_DROP, err.Error())
		return err
	}

	err = db.Remove(nil)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_DB_DROP, err.Error())
		return err
	}

	return nil
}

func (d *database) Table(name string) multidb.Table {
	col, err := d.getSession().Collection(nil, name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_NOTFOUND, err.Error())
	}

	return &table{Db: d, Name: name, Collection: col}
}

func (d *database) Graph(name string) multidb.Graph {
	return &graph{Db: d, Name: name}
}

func (d *database) ExecQuery(query string, params map[string]interface{}, res interface{}) (interface{}, error) {
	cursor, err := d.getSession().Query(nil, query, params)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_QUERY, err.Error())
		return nil, errors.New(multidb.EXCEPTION_QUERY)
	}
	defer cursor.Close()

	var docs []interface{}
	hasMore := true
	for {
		var doc interface{}
		_, err := cursor.ReadDocument(nil, &doc)
		if driver.IsNoMoreDocuments(err) {
			hasMore = false
		} else if err != nil {
			logrus.Errorf("%s: %s\n", multidb.EXCEPTION_QUERY, err.Error())
			return nil, errors.New(multidb.EXCEPTION_QUERY)
		}

		if !hasMore {
			break
		}
		docs = append(docs, doc)

	}

	if res != nil {
		multidb.ToStruct(docs, &res)
	}

	return docs, nil
}

// -----------------------------------------------
func (t *table) getSession() driver.Database {
	if t.Db.Session != nil {
		return t.Db.Session.(driver.Database)
	}
	return nil
}

func (t *table) Create() error {
	session := t.getSession()
	col, err := session.CreateCollection(nil, t.Name, nil)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_CREATE, err.Error())
		return err
	}
	t.Collection = col
	return nil
}

func (t *table) Exists() bool {
	session := t.getSession()
	exists, err := session.CollectionExists(nil, t.Name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_CONNECTION_FAILED, err.Error())
		return false
	}
	return exists
}

func (t *table) Drop() error {
	session := t.getSession()
	col, err := session.Collection(nil, t.Name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_DROP, err.Error())
		return err
	}

	err = col.Remove(nil)

	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_DROP, err.Error())
		return err
	}
	return nil
}

func (t *table) Insert(doc interface{}, res interface{}) (interface{}, error) {

	result, err := t.Collection.CreateDocument(nil, doc)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_INSERT_ERROR, err.Error())
		return nil, err
	}

	return result, nil
}

func (t *table) Find(id string, res interface{}) (interface{}, error) {
	doc, err := t.Collection.ReadDocument(nil, id, res)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_FIND_ERROR, err.Error())
		return nil, err
	}
	return doc, nil
}

func (t *table) Update(id string, doc interface{}, res interface{}) (interface{}, error) {

	result, err := t.Collection.UpdateDocument(nil, id, doc)

	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_UPDATE_ERROR, err.Error())
		return nil, err
	}

	if res != nil {
		multidb.ToStruct(result, &res)
	}

	return result, nil
}

func (t *table) Delete(id string, res interface{}) (interface{}, error) {
	result, err := t.Collection.RemoveDocument(nil, id)

	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_DELETE_ERROR, err.Error())
		return nil, err
	}

	if res != nil {
		multidb.ToStruct(result, &res)
	}

	return result, nil
}

//-----------------------------------------------
func (g *graph) getSession() driver.Database {
	if g.Db.Session != nil {
		return g.Db.Session.(driver.Database)
	}
	return nil
}

func (g *graph) Create() error {

	session := g.getSession()

	// A graph can contain additional vertex collections, defined in the set of orphan collections
	var options driver.CreateGraphOptions
	options.OrphanVertexCollections = []string{}
	options.EdgeDefinitions = []driver.EdgeDefinition{}

	// now it's possible to create a graph
	_, err := session.CreateGraph(nil, g.Name, &options)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_GRAPH_CREATE, err.Error())
		return err
	}

	return nil
}

func (g *graph) Relation(name string) multidb.Relation {
	col, err := g.getSession().Collection(nil, name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_NOTFOUND, err.Error())
	}

	return &relation{Name: name, Graph: g, Collection: col}
}

// -----------------------------------------------

func (r *relation) getSession() driver.Database {
	if r.Graph.Db.Session != nil {
		return r.Graph.Db.Session.(driver.Database)
	}
	return nil
}

func (r *relation) Create(from, to string, params interface{}) error {

	if !r.Exists() {
		session := r.getSession()
		graph, err := session.Graph(nil, r.Graph.Name)

		if err != nil {
			logrus.Errorf("%s: %s\n", multidb.EXCEPTION_RELATION_CREATE, err.Error())
			return err
		}

		collection, err := graph.CreateEdgeCollection(nil, r.Name, driver.VertexConstraints{
			From: []string{strings.Split(from, "/")[0]},
			To:   []string{strings.Split(to, "/")[0]},
		})

		if err != nil {
			logrus.Errorf("%s: %s\n", multidb.EXCEPTION_RELATION_CREATE, err.Error())
			return err
		}

		r.Collection = collection
	}

	return nil

}

func (r *relation) Exists() bool {
	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name, Collection: r.Collection})
	return tbl.Exists()
}

//
func (r *relation) Drop() error {
	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name})
	err := tbl.Drop()
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_RELATION_DELETE_ERROR, err.Error())
		return err
	}
	return nil
}

func (r *relation) Insert(data interface{}, res interface{}) (interface{}, error) {
	edge := data.(Edge)
	err := r.Create(edge.From, edge.To, nil)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_RELATION_INSERT_ERROR, err.Error())
		return nil, err
	}

	tbl := &table{Db: r.Graph.Db, Name: r.Name, Collection: r.Collection}
	relation, err := tbl.Insert(data, &res)

	if err != nil {
		return nil, errors.New(multidb.EXCEPTION_RELATION_INSERT_ERROR)
	}
	return relation, nil
}

func (r *relation) Find(id string, res interface{}) (interface{}, error) {
	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name})
	return tbl.Find(id, &res)
}

func (r *relation) Update(id string, doc interface{}, res interface{}) (interface{}, error) {
	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name})
	return tbl.Update(id, doc, &res)
}

func (r *relation) Delete(id string, res interface{}) (interface{}, error) {
	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name})
	return tbl.Delete(id, &res)
}
