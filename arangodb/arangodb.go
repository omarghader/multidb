package arangodb

//
// import (
// 	"errors"
// 	"fmt"
//
// 	arango "github.com/diegogub/aranGO"
// 	"github.com/sirupsen/logrus"
// )
//
// type arangodb struct {
// 	database
// }
//
// type arangoTable struct {
// 	table
// }
//
// type arangoGraph struct {
// 	graph
// }
//
// type arangoRelation struct {
// 	relation
// }
//
// func NewArangodb(options ConnectionOptions) Database {
// 	db := &arangodb{
// 		database{Name: options.DBName, Type: ARANGODB},
// 	}
//
// 	session, err := arango.Connect(fmt.Sprintf("http://%s:%s", options.Host, options.Port),
// 		options.Username, options.Password, true)
// 	if err != nil {
// 		logrus.Errorln(CONNECTION_FAILED)
// 		return nil
// 	}
//
// 	db.Session = session
// 	return db
// }
//
// func (db *arangodb) getSession() *arango.Session {
// 	if db.Session != nil {
// 		return db.Session.(*arango.Session)
// 	}
// 	return nil
// }
//
// func (db *arangodb) Create() error {
// 	session := db.getSession()
// 	return session.CreateDB(db.Name, nil)
// }
// func (db *arangodb) Exists() bool {
// 	session := db.getSession()
// 	availableDBs, err := session.AvailableDBs()
//
// 	if err != nil {
// 		return false
// 	}
//
// 	for _, availableDB := range availableDBs {
// 		if availableDB == db.Name {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// func (db *arangodb) Drop() error {
// 	return db.getSession().DropDB(db.Name)
//
// }
//
// func (db *arangodb) Table(name string) Table {
// 	return &arangoTable{table{Db: &db.database, Name: name}}
// }
//
// func (db *arangodb) Graph(name string, from, to []string) Graph {
// 	return &arangoGraph{graph{Db: &db.database, Name: name, fromType: from, toType: to}}
// }
//
// func (db *arangodb) ExecQuery(query string, params interface{}) ([]interface{}, error) {
// 	q := arango.NewQuery(query)
// 	cursor, err := db.getSession().DB(db.Name).Execute(q)
// 	if err != nil {
// 		return nil, errors.New(QUERY_EXCEPTION)
// 	}
//
// 	var res []interface{}
// 	err = cursor.FetchBatch(&res)
// 	if err != nil {
// 		return nil, errors.New(QUERY_EXCEPTION)
// 	}
//
// 	return res, nil
// }
//
// // -----------------------------------------------
// func (t *arangoTable) getSession() *arango.Session {
// 	if t.Db.Session != nil {
// 		return t.Db.Session.(*arango.Session)
// 	}
// 	return nil
// }
//
// func (t *arangoTable) Create() error {
// 	session := t.getSession()
// 	return session.DB(t.Db.Name).CreateCollection(arango.NewCollectionOptions(t.Name, true))
// }
//
// func (t *arangoTable) Exists() bool {
// 	session := t.getSession()
// 	return session.DB(t.Db.Name).ColExist(t.Name)
// }
//
// func (t *arangoTable) Drop() error {
// 	session := t.getSession()
// 	return session.DB(t.Db.Name).DropCollection(t.Name)
// }
// func (t *arangoTable) Insert(id string, doc interface{}) error {
// 	session := t.getSession()
// 	err := session.DB(t.Db.Name).Col(t.Name).Save(doc)
// 	if err != nil {
// 		return errors.New(INSERT_ERROR)
// 	}
// 	return nil
// }
//
// func (t *arangoTable) Find(id string) (interface{}, error) {
// 	session := t.getSession()
// 	var res interface{}
//
// 	err := session.DB(t.Db.Name).Col(t.Name).Get(id, &res)
// 	if err != nil {
// 		return nil, errors.New(FIND_ERROR)
// 	}
// 	return res, nil
// }
//
// func (t *arangoTable) Update(id string, doc interface{}) error {
// 	session := t.getSession()
// 	err := session.DB(t.Db.Name).Col(t.Name).Replace(id, doc)
// 	if err != nil {
// 		return errors.New(UPDATE_ERROR)
// 	}
// 	return nil
// }
// func (t *arangoTable) Delete(id string) error {
// 	session := t.getSession()
// 	err := session.DB(t.Db.Name).Col(t.Name).Delete(id)
// 	if err != nil {
// 		return errors.New(UPDATE_ERROR)
// 	}
// 	return nil
// }
//
// //-----------------------------------------------
// func (g *arangoGraph) getSession() *arango.Session {
// 	if g.Db.Session != nil {
// 		return g.Db.Session.(*arango.Session)
// 	}
// 	return nil
// }
//
// func (g *arangoGraph) Create() error {
//
// 	session := g.getSession()
// 	edges := arango.NewCollectionOptions(g.Name+"_edges", true)
// 	edges.IsEdge()
// 	err := session.DB(g.Db.Name).CreateCollection(edges)
// 	if err != nil {
// 		return err
// 	}
//
// 	edgeDefinition := arango.NewEdgeDefinition(g.Name+"_edges", g.fromType, g.toType)
// 	edgeDefinitions := []arango.EdgeDefinition{*edgeDefinition}
// 	_, err = session.DB(g.Db.Name).CreateGraph(g.Name, edgeDefinitions)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (g *arangoGraph) Relation(name string) Relation {
// 	return &arangoRelation{relation{Name: name, Graph: g.graph}}
// }
//
// // -----------------------------------------------
//
// func (r *arangoRelation) getSession() *arango.Session {
// 	if r.Graph.Db.Session != nil {
// 		return r.Graph.Db.Session.(*arango.Session)
// 	}
// 	return nil
// }
//
// func (r *arangoRelation) Create(from, to string, params interface{}) error {
// 	session := r.getSession()
// 	return session.DB(r.Graph.Db.Name).Col(r.Name).Relate(from, to, params)
// }
//
// func (r *arangoRelation) Exists() bool {
// 	// session := r.getSession()
// 	// return session.DB(r.Graph.Db.Name).Col(r.Graph.Name).
// 	return false
// }
//
// //
// func (r *arangoRelation) Drop() error {
// 	// session := t.getSession()
// 	// return session.DB(t.Db.Name).DropCollection(t.Name)
// 	return nil
// }
// func (r *arangoRelation) Insert(from, to string, params interface{}) error {
// 	// session := t.getSession()
// 	// err := session.DB(t.Db.Name).Col(t.Name).Save(doc)
// 	// if err != nil {
// 	// 	return errors.New(INSERT_ERROR)
// 	// }
// 	return nil
// }
//
// func (r *arangoRelation) Find(id string) (interface{}, error) {
// 	// session := t.getSession()
// 	var res interface{}
//
// 	// err := session.DB(t.Db.Name).Col(t.Name).Get(id, &res)
// 	// if err != nil {
// 	// 	return nil, errors.New(FIND_ERROR)
// 	// }
// 	return res, nil
// }
//
// func (r *arangoRelation) Update(id string, doc interface{}) error {
// 	// session := t.getSession()
// 	// err := session.DB(t.Db.Name).Col(t.Name).Replace(id, doc)
// 	// if err != nil {
// 	// 	return errors.New(UPDATE_ERROR)
// 	// }
// 	return nil
// }
// func (r *arangoRelation) Delete(id string) error {
// 	// session := t.getSession()
// 	// err := session.DB(t.Db.Name).Col(t.Name).Delete(id)
// 	// if err != nil {
// 	// 	return errors.New(UPDATE_ERROR)
// 	// }
// 	return nil
// }
