package multidb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/omarghader/multidb"
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
	if d.Session != nil {
		return d.Session.(driver.Database)
	}
	return nil
}

func (d *database) getClientSession() driver.Client {
	if d.ClientSession != nil {
		return d.ClientSession.(driver.Client)
	}
	return nil
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
	return &table{Db: d, Name: name}
}

func (d *database) Graph(name string, from, to []string) multidb.Graph {
	return &graph{Db: d, Name: name, fromType: from, toType: to}
}

func (d *database) ExecQuery(query string, params map[string]interface{}) ([]interface{}, error) {
	cursor, err := d.getSession().Query(nil, query, params)

	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_QUERY, err.Error())
		return nil, errors.New(multidb.EXCEPTION_QUERY)
	}
	defer cursor.Close()

	res := []interface{}{}

	for {
		var doc interface{}
		_, err := cursor.ReadDocument(nil, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logrus.Errorf("%s: %s\n", multidb.EXCEPTION_QUERY, err.Error())
			return nil, errors.New(multidb.EXCEPTION_QUERY)
		}

		res = append(res, doc)
	}

	return res, nil
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
	_, err := session.CreateCollection(nil, t.Name, nil)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_CREATE, err.Error())
		return err
	}
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

func (t *table) Insert(id string, doc interface{}) (interface{}, error) {
	session := t.getSession()
	col, err := session.Collection(nil, t.Name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_INSERT_ERROR, err.Error())
		return nil, err
	}

	res, err := col.CreateDocument(nil, doc)

	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_INSERT_ERROR, err.Error())
		return nil, err
	}

	return res, nil
}

func (t *table) Find(id string) (interface{}, error) {
	session := t.getSession()
	var res map[string]interface{}

	col, err := session.Collection(nil, t.Name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_FIND_ERROR, err.Error())
		return nil, err
	}

	_, err = col.ReadDocument(nil, id, &res)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_FIND_ERROR, err.Error())
		return nil, err
	}
	return res, nil
}

func (t *table) Update(id string, doc interface{}) (interface{}, error) {
	session := t.getSession()
	col, err := session.Collection(nil, t.Name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_UPDATE_ERROR, err.Error())
		return nil, err
	}

	res, err := col.UpdateDocument(nil, id, doc)

	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_UPDATE_ERROR, err.Error())
		return nil, err
	}

	return res, nil
}

func (t *table) Delete(id string) (interface{}, error) {
	session := t.getSession()
	col, err := session.Collection(nil, t.Name)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_DELETE_ERROR, err.Error())
		return nil, err
	}

	res, err := col.RemoveDocument(nil, id)

	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_TABLE_DELETE_ERROR, err.Error())
		return nil, err
	}
	return res, nil
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
	return &relation{Name: name, Graph: g}
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

		_, err = graph.CreateEdgeCollection(nil, r.Name, driver.VertexConstraints{
			From: []string{strings.Split(from, "/")[0]},
			To:   []string{strings.Split(to, "/")[0]},
		})

		if err != nil {
			logrus.Errorf("%s: %s\n", multidb.EXCEPTION_RELATION_CREATE, err.Error())
			return err
		}

	}
	return nil

}

func (r *relation) Exists() bool {
	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name})
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

func (r *relation) Insert(from, to string, props interface{}) (interface{}, error) {
	err := r.Create(from, to, nil)
	if err != nil {
		logrus.Errorf("%s: %s\n", multidb.EXCEPTION_RELATION_INSERT_ERROR, err.Error())
		return nil, err
	}

	marshalledProps, err := json.Marshal(props)
	if err != nil {
		logrus.Errorf("%s: %s\n", "Cannot parse relation props", err.Error())
		return nil, err
	}

	var relProps map[string]interface{}

	err = json.Unmarshal(marshalledProps, &relProps)
	if err != nil {
		logrus.Errorf("%s: %s\n", "Cannot parse relation props", err.Error())
		return nil, err
	}

	relProps["_from"] = from
	relProps["_to"] = to

	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name})
	relation, err := tbl.Insert("", relProps)

	if err != nil {
		return nil, errors.New(multidb.EXCEPTION_RELATION_INSERT_ERROR)
	}
	return relation, nil
}

func (r *relation) Find(id string) (interface{}, error) {
	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name})
	return tbl.Find(id)
}

func (r *relation) Update(id string, doc interface{}) (interface{}, error) {
	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name})
	return tbl.Update(id, doc)
}

func (r *relation) Delete(id string) (interface{}, error) {
	tbl := multidb.Table(&table{Db: r.Graph.Db, Name: r.Name})
	return tbl.Delete(id)
}
