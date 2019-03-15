package multidb

import (
	"fmt"
	"testing"

	arango "github.com/diegogub/aranGO"
)

// func TestArangoDatabase(t *testing.T) {
// 	db := NewArangodb(ConnectionOptions{
// 		Host:     "localhost",
// 		Port:     "",
// 		Username: "root",
// 		Password: "root",
// 		DBName:   "test",
// 	})
//
// 	if db != nil {
// 		t.Errorf("db should be nil")
// 	}
//
// 	db = NewArangodb(ConnectionOptions{
// 		Host:     "localhost",
// 		Port:     "8529",
// 		Username: "root",
// 		Password: "root",
// 		DBName:   "test",
// 	})
//
// 	if db == nil {
// 		t.Errorf(CONNECTION_FAILED)
// 	}
//
// 	if db.Exists() {
// 		err := db.Drop()
// 		if err != nil {
// 			t.Errorf("Cannot drop database")
// 		}
// 	}
//
// 	err := db.Create()
// 	if err != nil {
// 		t.Errorf("Cannot create database")
// 	}
//
// }
//
// func TestArangoTable(t *testing.T) {
// 	db := NewArangodb(ConnectionOptions{
// 		Host:     "localhost",
// 		Port:     "8529",
// 		Username: "root",
// 		Password: "root",
// 		DBName:   "test",
// 	})
//
// 	tbl := db.Table("table_test")
// 	if tbl.Exists() {
// 		err := tbl.Drop()
// 		if err != nil {
// 			t.Errorf("Cannot drop table")
// 		}
// 	}
//
// 	err := tbl.Create()
// 	if err != nil {
// 		t.Errorf("Cannot create table")
// 	}
//
// 	doc := map[string]string{"_key": "doc1", "name": "document1"}
// 	err = tbl.Insert("doc1", doc)
// 	if err != nil {
// 		t.Errorf("Cannot insert table")
// 	}
//
// 	doc1, err := tbl.Find("doc1")
// 	if err != nil {
// 		t.Errorf("Cannot find table")
// 	}
//
// 	if doc1.(map[string]interface{})["name"] != "document1" {
// 		t.Errorf("Cannot find object from table")
// 	}
//
// 	doc["name"] = "document2"
// 	err = tbl.Update("doc1", doc)
// 	if err != nil {
// 		t.Errorf("Cannot update table")
// 	}
//
// 	doc2, err := tbl.Find("doc1")
// 	if err != nil {
// 		t.Errorf("Cannot find table")
// 	}
//
// 	if doc2.(map[string]interface{})["name"] != "document2" {
// 		t.Errorf("Cannot updae object from table")
// 	}
//
// 	err = tbl.Delete("doc1")
// 	if err != nil {
// 		t.Errorf("Cannot find table")
// 	}
//
// 	tbl.Drop()
// }

type DocTest struct {
	arango.Document // Must include arango Document in every struct you want to save id, key, rev after saving it
	Name            string
	Age             int
	Likes           []string
}

func TestArangoGraph(t *testing.T) {
	// 	db := NewArangodb(ConnectionOptions{
	// 		Host:     "localhost",
	// 		Port:     "8529",
	// 		Username: "root",
	// 		Password: "root",
	// 		DBName:   "test",
	// 	})
	//
	// 	db.Table("test1").Create()
	// 	db.Table("test2").Create()
	// 	// err := db.Graph("graphy", []string{"test1"}, []string{"test2"}).Create()
	// 	// fmt.Println(err)
	//
	// 	var d1, d2 DocTest
	// 	d1.Name = "Diego"
	// 	d1.Age = 22
	// 	d1.Likes = []string{"arangodb", "golang", "linux"}
	//
	// 	d2.Name = "Facundo"
	// 	d2.Age = 25
	// 	d2.Likes = []string{"php", "linux", "python"}
	//
	// 	err := db.Table("test1").Insert("", d1)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	err = db.Table("test2").Insert("", d2)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	var graphName = "graphiti"
	// 	graph := db.Graph(graphName, []string{"test1"}, []string{"test2"})
	// 	fmt.Println(graph.Create())
	//
	// 	// rel := graph.Relation(graphName + "edges")
	// 	// // Relate documents
	// 	// err = rel.Insert(d1.Id, d2.Id, map[string]interface{}{"is": "friend"})
	// 	//
	// 	// fmt.Println(err)
}

func TestX(t *testing.T) {
	//change false to true if you want to see every http request
	//Connect(host, user, password string, log bool) (*Session, error) {
	s, err := arango.Connect("http://localhost:8529", "root", "root", true)
	if err != nil {
		panic(err)
	}

	// CreateDB(name string,users []User) error
	s.CreateDB("test", nil)

	// create Collections test if exist
	if !s.DB("test").ColExist("docs1") {
		// CollectionOptions has much more options, here we just define name , sync
		docs1 := arango.NewCollectionOptions("docs1", true)
		s.DB("test").CreateCollection(docs1)
	}

	if !s.DB("test").ColExist("docs2") {
		docs2 := arango.NewCollectionOptions("docs2", true)
		s.DB("test").CreateCollection(docs2)
	}

	if !s.DB("test").ColExist("ed") {
		edges := arango.NewCollectionOptions("ed", true)
		edges.IsEdge() // set to Edge
		s.DB("test").CreateCollection(edges)
	}
	// Create and Relate documents
	var d1, d2 DocTest
	d1.Name = "Diego"
	d1.Age = 22
	d1.Likes = []string{"arangodb", "golang", "linux"}

	d2.Name = "Facundo"
	d2.Age = 25
	d2.Likes = []string{"php", "linux", "python"}

	err = s.DB("test").Col("docs1").Save(&d1)
	err = s.DB("test").Col("docs1").Save(&d2)
	if err != nil {
		panic(err)
	}

	ed := arango.NewEdgeDefinition("ed1", []string{"docs1"}, []string{"docs2"})
	fmt.Println(s.DB("test").CreateGraph("graphii", []arango.EdgeDefinition{*ed}))

	// Relate documents
	// s.DB("test").Col("ed").Relate(d1.Id, d2.Id, map[string]interface{}{"is": "friend"})

}
